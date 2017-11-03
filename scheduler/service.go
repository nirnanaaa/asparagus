package scheduler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nirnanaaa/asparagus/etcd"

	"github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"github.com/gorhill/cronexpr"
	metrics "github.com/rcrowley/go-metrics"
)

const (
	metricNameExecutions = "cronjobExecutions"
)

type ServiceConfig struct {
	TickDuration int64  `json:"SchedulerTickDuration"`
	Enabled      bool   `json:"Enabled"`
	SecretKey    string `json:"SecretKey"`
}

// Service Returns a config
type Service struct {
	ServiceConfig ServiceConfig
	Config        *Config
	Logger        *logrus.Logger
	ETCD          *etcd.Service
	Cancel        chan struct{}
	Ticker        *time.Ticker
}

func NewService(config *Config) *Service {
	return &Service{
		Config: config,
		Cancel: make(chan struct{}),
		Logger: logrus.New(),
	}
}
func (s *Service) Printf(format string, v ...interface{}) {
	s.Logger.Infof(format, v)
}

// Start starts the actual service
func (s *Service) Start() error {
	if err := s.InitConfig(); err != nil {
		return err
	}
	s.Ticker = time.NewTicker(time.Millisecond * time.Duration(s.ServiceConfig.TickDuration))
	tasks, err := s.getJobs()
	if err != nil {
		return err
	}
	for {
		select {
		case <-s.Ticker.C:
			if !s.ServiceConfig.Enabled {
				s.Logger.Infof("cron daemon is disabled. Please make sure to configure it properly using etcd")
			}
			m := metrics.GetOrRegisterCounter(metricNameExecutions, metrics.DefaultRegistry)
			m.Inc(1)
			s.Tick(tasks)
		case <-s.Cancel:
			return nil
		}
	}
}

// Stop cancels any executions immediately
func (s *Service) Stop() error {
	s.Cancel <- struct{}{}
	<-s.Cancel
	return nil
}

func (s *Service) InitConfig() error {
	watcherOpts := client.WatcherOptions{AfterIndex: 0}
	go func() {
		w := s.ETCD.KAPI.Watcher(s.Config.FolderHotConfig, &watcherOpts)
		for {
			r, err := w.Next(context.Background())
			if err != nil {
				log.Fatal("Error occurred", err)
			}
			if r.Action != "set" {
				continue
			}
			conf, err := s.ETCD.KAPI.Get(context.Background(), s.Config.FolderHotConfig, nil)
			if err != nil {
				return
			}
			value := conf.Node.Value
			var cfg ServiceConfig
			if err := json.Unmarshal([]byte(value), &cfg); err != nil {
				continue
			}
			s.ServiceConfig = cfg
		}
	}()
	conf, err := s.ETCD.KAPI.Get(context.Background(), s.Config.FolderHotConfig, nil)
	if err != nil {
		return s.DefaultServiceConfig()
	}
	value := conf.Node.Value
	var cfg ServiceConfig
	if err := json.Unmarshal([]byte(value), &cfg); err != nil {
		return err
	}
	s.ServiceConfig = cfg
	return nil
}

func (s *Service) DefaultServiceConfig() error {
	cfg := ServiceConfig{
		TickDuration: 30000,
		Enabled:      true,
	}
	scfg, err := json.Marshal(&cfg)
	if err != nil {
		return err
	}
	if _, err := s.ETCD.KAPI.Set(context.Background(), s.Config.FolderHotConfig, string(scfg), nil); err != nil {
		return err
	}
	s.ServiceConfig = cfg
	return nil
}

func (s *Service) getJobs() (*Tasks, error) {
	runner := NewTasksWithKeys(s.ETCD.KAPI)
	runner.Watch()
	if s.Config.LogTasksDetection {
		runner.DebugTasks(s.Logger)
	}
	if err := runner.Load(); err != nil {
		return nil, err
	}
	return runner, nil
}
func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
func (s *Service) Tick(tasks *Tasks) error {
	for _, task := range tasks.Archive {
		if task.Running {
			continue
		}
		expr, err := cronexpr.Parse(task.Expression)
		if err != nil {
			continue
		}
		expEnd := expr.Next(task.LastRunAt)
		if time.Now().Before(expEnd) {
			continue
		}
		if task.AfterTask != "" {
			dependsTask, err := tasks.GetTask(task.AfterTask)
			if err != nil {
				continue
			}
			if dependsTask.Running {
				continue
			}
			if dependsTask.LastRunAt.Before(Bod(time.Now())) {
				// depending task hasn't executed, yet
				continue
			}
		}
		task.SetRunning(s.ETCD.KAPI, true)
		executor := NewExecutor([]byte(s.ServiceConfig.SecretKey), s.Logger)
		executor.HTTPConfig = s.Config.HTTPConfig
		go executor.FromTask(task)
		task.Done(s.ETCD.KAPI)
		tasks.Archive[task.Name] = task
	}
	return nil
}
