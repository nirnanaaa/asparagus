package scheduler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	metrics "github.com/rcrowley/go-metrics"
)

const (
	metricNameExecutionSuccessful   = "cronjobExecutionSuccessful"
	metricNameExecutionSuccessMeter = "cronjobExecutionSuccessMeter"
	metricNameExecutionTiming       = "cronjobExecutionTiming"
	metricNameExecutionFailed       = "cronjobExecutionFailed"
	metricNameExecutionError        = "cronjobExecutionErrorMeter"
	metricNameExecutionPayloadSize  = "cronjobExecutionPayloadSize"
)

// Executor executes jobs
type Executor struct {
	secretKey      []byte
	Logger         *logrus.Logger
	ErrorCounter   metrics.Counter
	SuccessCounter metrics.Counter
	RequestTiming  metrics.Timer
	HTTPConfig     HTTPConfig
}

// NewExecutor returns a new executionService
func NewExecutor(secret []byte, logger *logrus.Logger) Executor {
	errC := metrics.GetOrRegisterCounter(metricNameExecutionError, metrics.DefaultRegistry)
	okC := metrics.GetOrRegisterCounter(metricNameExecutionSuccessful, metrics.DefaultRegistry)
	rTime := metrics.GetOrRegisterTimer(metricNameExecutionTiming, metrics.DefaultRegistry)
	return Executor{
		ErrorCounter:   errC,
		SuccessCounter: okC,
		RequestTiming:  rTime,
		secretKey:      secret,
		Logger:         logger,
	}
}

// FromTask runs a task definition
func (e *Executor) FromTask(t TaskDefinition) error {
	_, err := e.request(t.URI, t.Name)
	return err
}

func (e *Executor) request(url, name string) (response *http.Response, erro error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		Issuer:    "/",
		Subject:   "0",
	}
	e.RequestTiming.Time(func() {

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := token.SignedString(e.secretKey)
		if err != nil {
			e.ErrorCounter.Inc(1)
			erro = err
			return
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			e.ErrorCounter.Inc(1)
			erro = err
			return
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ss))
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			e.ErrorCounter.Inc(1)
			erro = err
			return
		}
		e.logHTTPStatusCode(name, req, resp)
		e.logHTTPResponseBody(name, req, resp)
		if resp.StatusCode >= 400 {
			erro = fmt.Errorf("Request failed.")
		}
		erro = nil
		response = resp
	})
	return
}

func (e *Executor) logHTTPResponseBody(name string, req *http.Request, resp *http.Response) {
	if !e.HTTPConfig.LogResponseBody {
		return
	}
	defer resp.Body.Close()
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		e.Logger.Errorf("%s: Response Body dumping errored: \n %s", name, err.Error())
		return
	}
	e.Logger.Infof("%s: Response Body returned: \n %q", name, dump)
}

func (e *Executor) logHTTPStatusCode(name string, req *http.Request, resp *http.Response) {
	if !e.HTTPConfig.LogResponseStatus {
		return
	}
	if resp.StatusCode >= 400 {
		e.ErrorCounter.Inc(1)
		e.Logger.Warnf("%s: request to %s failed with response code: %d", name, req.URL.String(), resp.StatusCode)
		return
	}
	e.SuccessCounter.Inc(1)
	e.Logger.Debugf("%s: request to %s succeeded with response code: %d", name, req.URL.String(), resp.StatusCode)
}
