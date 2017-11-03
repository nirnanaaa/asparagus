package metric

import (
	"github.com/Sirupsen/logrus"
	"github.com/nirnanaaa/asparagus/metric/adapters"
	"github.com/nirnanaaa/asparagus/metric/adapters/cloudwatch"
	"github.com/nirnanaaa/asparagus/metric/adapters/console"
	"github.com/nirnanaaa/asparagus/metric/adapters/slack"
)

// Service Returns a config
type Service struct {
	Config    *Config
	Logger    *logrus.Logger
	Reporters []adapters.Reporter
}

// NewService creates a new service adapter
func NewService(config *Config) *Service {
	reporters := []adapters.Reporter{
		cloudwatch.NewReporter(config.Cloudwatch),
		slack.NewReporter(config.Slack),
		console.NewReporter(config.Console),
	}
	return &Service{
		Config:    config,
		Reporters: reporters,
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

	for _, reporter := range s.Reporters {
		if err := reporter.StartReporting(); err != nil {
			return err
		}
	}
	s.Logger.Info("Started Metrics Service")
	return nil
}

// Stop stops the metrics service
func (s *Service) Stop() error {
	return nil
}
