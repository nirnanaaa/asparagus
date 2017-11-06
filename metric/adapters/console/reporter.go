package console

import (
	"log"
	"os"
	"time"

	"github.com/nirnanaaa/asparagus/metric/adapters"
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

// SupportsLogging returns if the adapter supports logging
func (r *Reporter) SupportsLogging() bool {
	return false
}

// LogResult logs a log line if SupportsLogging is enabled
func (r *Reporter) LogResult(evt *adapters.LogEvent) error {
	return nil
}
