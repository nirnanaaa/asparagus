package cloudwatch

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cw "github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/nirnanaaa/asparagus/metric/adapters/cloudwatch/config"
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
	metricsRegistry := metrics.DefaultRegistry
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	}))
	registrator := &config.Config{
		Client:                cw.New(sess),
		Namespace:             r.Config.Namespace,
		Filter:                &config.NoFilter{},
		ResetCountersOnReport: r.Config.ResetCountersAfterReport,
		ReportingInterval:     time.Duration(r.Config.ReportingInterval),
		StaticDimensions:      r.Config.StaticDimensions,
	}
	go Cloudwatch(metricsRegistry, registrator)
	return nil
}

// StopReporting stops the recording process
func (r *Reporter) StopReporting() error {
	return nil
}
