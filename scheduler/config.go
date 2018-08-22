package scheduler

import (
	"time"

	"github.com/nirnanaaa/asparagus/scheduler/provider/crontab"
	"github.com/nirnanaaa/asparagus/scheduler/provider/etcd"
	"github.com/nirnanaaa/asparagus/scheduler/provider/http"
	"github.com/nirnanaaa/asparagus/scheduler/provider/local"
	"github.com/nirnanaaa/asparagus/scheduler/provider/sql"
	"github.com/nirnanaaa/asparagus/toml"
)

const (
	//DefaultTickDuration is default state for metrics
	DefaultTickDuration = 1
)

// Config represents the meta configuration.
type Config struct {
	Enabled           bool           `toml:"enabled"`
	TickDuration      toml.Duration  `toml:"tick-duration"`
	LogTasksDetection bool           `toml:"log-task-detection"`
	HTTPExecutor      http.Config    `toml:"executor-http"`
	LocalExecutor     local.Config   `toml:"executor-local"`
	CrontabSource     crontab.Config `toml:"provider-crontab"`
	ETCDSource        etcd.Config    `toml:"provider-etcd"`
	SQLSource         sql.Config     `toml:"provider-sql"`
	NumWorkers        int            `toml:"num-workers"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		Enabled:       true,
		TickDuration:  toml.Duration(DefaultTickDuration * time.Second),
		HTTPExecutor:  http.NewConfig(),
		LocalExecutor: local.NewConfig(),
		CrontabSource: crontab.NewConfig(),
		ETCDSource:    etcd.NewConfig(),
		SQLSource:     sql.NewConfig(),
		NumWorkers:    10,
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	return nil
}
