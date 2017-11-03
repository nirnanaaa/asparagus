package metric

import (
	"log"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	cw "github.com/nirnanaaa/asparagus/metric/adapters/cloudwatch"
	"github.com/nirnanaaa/asparagus/metric/adapters/cloudwatch/config"
	metrics "github.com/rcrowley/go-metrics"
)

// Service Returns a config
type Service struct {
	Config *Config
	Logger *logrus.Logger
}

// NewService creates a new service adapter
func NewService(config *Config) *Service {
	return &Service{
		Config: config,
	}
}

// Printf formats a message ot the default logger
func (s *Service) Printf(format string, v ...interface{}) {
	s.Logger.Infof(format, v)
}

// Start starts the metrics reporter
func (s *Service) Start() error {
	if !s.Config.Enabled {
		return nil
	}
	s.Logger.Info("Starting Metrics Service")
	go metrics.Log(
		metrics.DefaultRegistry,
		21600*time.Second,
		log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	metricsRegistry := metrics.DefaultRegistry
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	}))
	registrator := &config.Config{
		Client:                cloudwatch.New(sess),
		Namespace:             "ASPARAGUS",
		Filter:                &config.NoFilter{},
		ResetCountersOnReport: true,
		ReportingInterval:     10 * time.Second,
		StaticDimensions: map[string]string{
			"Environment": os.Getenv("ENV"),
		},
	}
	go cw.Cloudwatch(metricsRegistry, registrator)
	return nil
}

// Stop stops the metrics service
func (s *Service) Stop() error {
	return nil
}
