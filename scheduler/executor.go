package scheduler

import (
	"fmt"
	"net/http"
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
	secretKey []byte
	Logger    *logrus.Logger
}

// NewExecutor returns a new executionService
func NewExecutor(secret []byte, logger *logrus.Logger) Executor {
	return Executor{
		secretKey: secret,
		Logger:    logger,
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
		Issuer:    "https://logx.miha-bodytec.com/",
		Subject:   "0",
	}
	errC := metrics.GetOrRegisterCounter(metricNameExecutionFailed, metrics.DefaultRegistry)
	errM := metrics.GetOrRegisterMeter(metricNameExecutionError, metrics.DefaultRegistry)
	t := metrics.GetOrRegisterTimer(metricNameExecutionTiming, metrics.DefaultRegistry)
	t.Time(func() {

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := token.SignedString(e.secretKey)
		if err != nil {
			errC.Inc(1)
			errM.Mark(1)
			erro = err
			return
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			errC.Inc(1)
			errM.Mark(1)
			erro = err
			return
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ss))
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errC.Inc(1)
			errM.Mark(1)
			erro = err
			return
		}
		if resp.StatusCode >= 400 {
			errC.Inc(1)
			errM.Mark(1)
			e.Logger.Warnf("%s: request to %s failed with response code: %d", name, url, resp.StatusCode)
			erro = fmt.Errorf("response returned a bad response code: %d", resp.StatusCode)
			return
		}

		okC := metrics.GetOrRegisterCounter(metricNameExecutionSuccessful, metrics.DefaultRegistry)
		okM := metrics.GetOrRegisterMeter(metricNameExecutionSuccessMeter, metrics.DefaultRegistry)
		okC.Inc(1)
		okM.Mark(1)
		e.Logger.Debugf("%s: request to %s succeeded with response code: %d", name, url, resp.StatusCode)
		erro = nil
		response = resp
	})
	return
}
