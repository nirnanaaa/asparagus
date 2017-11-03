package console

import (
	"log"
	"os"
	"time"

	metrics "github.com/rcrowley/go-metrics"
)

// Reporter type
type Reporter struct {
	Config Config
}

// NewReporter returns a new reporter
func NewReporter(config Config) *Reporter {
	return &Reporter{
		config,
	}
}

// StartReporting starts the recording process
func (r *Reporter) StartReporting() error {
	if !r.Config.Enabled {
		return nil
	}
	go metrics.Log(
		metrics.DefaultRegistry,
		time.Duration(r.Config.ReportingInterval),
		log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
	return nil
}

// StopReporting stops the recording process
func (r *Reporter) StopReporting() error {
	return nil
}
