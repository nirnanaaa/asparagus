package scheduler

import (
	"time"

	"github.com/nirnanaaa/asparagus/scheduler/provider"
	"github.com/nirnanaaa/asparagus/scheduler/provider/http"
	"github.com/nirnanaaa/asparagus/scheduler/provider/local"

	"github.com/Sirupsen/logrus"
	"github.com/gorhill/cronexpr"
	metrics "github.com/rcrowley/go-metrics"
)

const (
	metricNameExecutions = "cronjobExecutions"
)

// Service Returns a config
type Service struct {
	Config             *Config
	Logger             *logrus.Logger
	Cancel             chan struct{}
	Ticker             *time.Ticker
	ExecutionProviders map[string]provider.ExecutionProvider
	JobDispatcher      *Dispatcher
}

// NewService creates a new service
func NewService(config *Config, logger *logrus.Logger) *Service {
	target := map[string]provider.ExecutionProvider{
		"http":  http.NewExecutionProvider(config.HTTPExecutor, logger),
		"local": local.NewExecutionProvider(config.LocalExecutor, logger),
	}
	return &Service{
		Config:             config,
		Cancel:             make(chan struct{}),
		Logger:             logger,
		ExecutionProviders: target,
	}
}

// Start starts the actual service
func (s *Service) Start() error {
	if !s.Config.Enabled {
		s.Logger.Warn("Scheduler is disabled in config.")
		return nil
	}
	dispatcher := StartDispatcher(s.Config.NumWorkers, s.ExecutionProviders)
	s.JobDispatcher = dispatcher
	s.Ticker = time.NewTicker(time.Duration(s.Config.TickDuration))
	tasks, err := s.getJobs()
	if err != nil {
		return err
	}
	ticks := metrics.GetOrRegisterCounter(metricNameExecutions, metrics.DefaultRegistry)
	for {
		select {
		case <-s.Ticker.C:
			ticks.Inc(1)
			s.Tick(tasks)
		case <-s.Cancel:
			return nil
		}
	}
}

// Stop cancels any executions immediately
func (s *Service) Stop() error {
	s.Cancel <- struct{}{}
	// <-s.Cancel
	return nil
}

// RegisterExecutionProvider registers or overwrites an execution provider
func (s *Service) RegisterExecutionProvider(name string, provider provider.ExecutionProvider) {
	s.ExecutionProviders[name] = provider
}

func (s *Service) getJobs() (*Tasks, error) {
	runner := NewSourceConfig(s.Config)

	if err := runner.Load(); err != nil {
		return nil, err
	}
	if s.Config.LogTasksDetection {
		s.Logger.Infof("Task detection is been used")
		runner.DebugTasks(s.Logger)
	}
	return runner, nil
}

// Bod returns the beginning of the day
func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// Tick runs for a set interval
func (s *Service) Tick(tasks *Tasks) error {
	for _, task := range tasks.Tasks {
		if task.Running {
			continue
		}
		expr, err := cronexpr.Parse(task.Expression)
		if err != nil {
			s.Logger.WithField("Task", task.Name).WithError(err).Debug("parsing expression failed")
			continue
		}
		expEnd := expr.Next(task.LastRunAt)
		if time.Now().Before(expEnd) {
			s.Logger.WithField("Task", task.Name).Debug("cronjob shouldn't been executed, yet")
			continue
		}
		if task.AfterTask != "" {
			s.Logger.WithField("Task", task.Name).Debug("has dependencies")
			dependsTask, err := tasks.GetTask(task.AfterTask)
			if err != nil {
				continue
			}
			if dependsTask.Running {
				s.Logger.WithField("Task", task.Name).Debug("dependencies are currently running")
				continue
			}
			if dependsTask.LastRunAt.Before(Bod(time.Now())) {
				s.Logger.WithField("Task", task.Name).Debug("dependencies didn't run today, yet")
				continue
			}
		}
		s.Logger.WithField("Task", task.Name).Info("Queueing job")
		s.JobDispatcher.WorkQueue <- task
	}
	return nil
}
