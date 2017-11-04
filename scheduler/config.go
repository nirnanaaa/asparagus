package scheduler

import (
	"time"

	"github.com/nirnanaaa/asparagus/scheduler/provider/crontab"
	"github.com/nirnanaaa/asparagus/scheduler/provider/http"
	"github.com/nirnanaaa/asparagus/toml"
)

const (
	//DefaultTickDuration is default state for metrics
	DefaultTickDuration = 1

	//DefaultAPIFallbackDomain defines a default fallback
	DefaultAPIFallbackDomain = "example.com"

	// DefaultFolderCronjobs defines a safe-default
	DefaultFolderCronjobs = "/cron/Jobs"

	// DefaultHotConfigPath defines a key that is used for live re-configuration
	DefaultHotConfigPath = "/cron/Config"
)

// HTTPConfig is used to describe all related HTTP configuration
type HTTPConfig struct {
	LogResponseBody   bool `toml:"log-response-body"`
	LogResponseStatus bool `toml:"log-response-status"`
}

// Config represents the meta configuration.
type Config struct {
	Enabled           bool           `toml:"enabled"`
	TickDuration      toml.Duration  `toml:"tick-duration"`
	APIFallbackDomain string         `toml:"api-fallback-domain"`
	FolderCronjobs    string         `toml:"cronjob-base-folder"`
	FolderHotConfig   string         `toml:"hot-config-path"`
	HTTPConfig        HTTPConfig     `toml:"http"`
	LogTasksDetection bool           `toml:"log-task-detection"`
	HTTPExecutor      http.Config    `toml:"executor-http"`
	CrontabSource     crontab.Config `toml:"provider-crontab"`
	NumWorkers        int            `toml:"num-workers"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		Enabled:           true,
		TickDuration:      toml.Duration(DefaultTickDuration * time.Second),
		APIFallbackDomain: DefaultAPIFallbackDomain,
		FolderCronjobs:    DefaultFolderCronjobs,
		FolderHotConfig:   DefaultHotConfigPath,
		HTTPExecutor:      http.NewConfig(),
		CrontabSource:     crontab.NewConfig(),
		NumWorkers:        10,
	}
}

func (c *Config) Validate() error {
	return nil
}
