package http

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/nirnanaaa/asparagus/scheduler/provider"
	metrics "github.com/rcrowley/go-metrics"
)

// Executor executes jobs
type Executor struct {
	Logger         *logrus.Logger
	ErrorCounter   metrics.Counter
	SuccessCounter metrics.Counter
	RequestTiming  metrics.Timer
	HTTPConfig     Config
}

// NewExecutor returns a new executionService
func NewExecutor(logger *logrus.Logger) Executor {
	errC := metrics.GetOrRegisterCounter(metricNameExecutionError, metrics.DefaultRegistry)
	okC := metrics.GetOrRegisterCounter(metricNameExecutionSuccessful, metrics.DefaultRegistry)
	rTime := metrics.GetOrRegisterTimer(metricNameExecutionTiming, metrics.DefaultRegistry)
	return Executor{
		ErrorCounter:   errC,
		SuccessCounter: okC,
		RequestTiming:  rTime,
		Logger:         logger,
	}
}

// FromTask runs a task definition
func (e *Executor) FromTask(t ExecutionData, res *provider.Response) error {
	if _, err := e.Request(t.URL, t.Method, res); err != nil {
		return err
	}
	return nil
}

// Request performs an actual http request
func (e *Executor) Request(url, method string, resq *provider.Response) (response *http.Response, erro error) {
	e.RequestTiming.Time(func() {
		req, err := http.NewRequest(strings.ToUpper(method), url, nil)
		if err != nil {
			e.ErrorCounter.Inc(1)
			erro = err
			return
		}
		if e.HTTPConfig.SignJWT {
			claims := jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Duration(e.HTTPConfig.JWTExpires)).Unix(),
				Issuer:    e.HTTPConfig.Issuer,
				Subject:   e.HTTPConfig.JWTSubject,
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			ss, err := token.SignedString([]byte(e.HTTPConfig.JWTSecret))
			if err != nil {
				e.ErrorCounter.Inc(1)
				erro = err
				return
			}
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ss))
		}
		e.logHTTPRequest(method, req)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			e.ErrorCounter.Inc(1)
			erro = err
			return
		}
		e.logHTTPStatusCode(method, req, resp, resq)
		e.logHTTPResponseBody(method, req, resp, resq)
		if resp.StatusCode >= 400 {
			erro = fmt.Errorf("request failed with status code %d", resp.StatusCode)
			return
		}
		erro = nil
		response = resp
	})
	return
}

func (e *Executor) logHTTPRequest(name string, req *http.Request) {
	if !e.HTTPConfig.DebugResponse {
		return
	}
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		e.Logger.
			Errorf("%s: Request dumping errored: \n%s", name, err.Error())
		return
	}
	e.Logger.Infof("%s: Request returned: \n%s", name, dump)
}
func (e *Executor) logHTTPResponseBody(name string, req *http.Request, resp *http.Response, resq *provider.Response) {
	defer resp.Body.Close()
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		if e.HTTPConfig.DebugResponse {
			e.Logger.
				Errorf("%s: Response Body dumping errored: \n%s", name, err.Error())
		}
		return
	}
	resq.Response = dump
	if e.HTTPConfig.DebugResponse {
		e.Logger.Infof("%s: Response Body returned: \n%s", name, dump)
	}
}

func (e *Executor) logHTTPStatusCode(name string, req *http.Request, resp *http.Response, resq *provider.Response) {
	resq.StatusCode = resp.StatusCode
	if resp.StatusCode >= 400 {
		e.ErrorCounter.Inc(1)
		e.Logger.Warnf("%s: request to %s failed with response code: %d", name, req.URL.String(), resp.StatusCode)
		return
	}
	e.SuccessCounter.Inc(1)
	if e.HTTPConfig.LogHTTPStatus {
		e.Logger.Debugf("%s: request to %s succeeded with response code: %d", name, req.URL.String(), resp.StatusCode)
	}
}
