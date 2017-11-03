package scheduler

const (
	//DefaultTickDuration is default state for metrics
	DefaultTickDuration = 1000

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
	TickDuration      int64      `toml:"tick-duration-ms"`
	APIFallbackDomain string     `toml:"api-fallback-domain"`
	FolderCronjobs    string     `toml:"cronjob-base-folder"`
	FolderHotConfig   string     `toml:"hot-config-path"`
	HTTPConfig        HTTPConfig `toml:"http"`
	LogTasksDetection bool       `toml:"log-task-detection"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		TickDuration:      DefaultTickDuration,
		APIFallbackDomain: DefaultAPIFallbackDomain,
		FolderCronjobs:    DefaultFolderCronjobs,
		FolderHotConfig:   DefaultHotConfigPath,
	}
}

func (c *Config) Validate() error {
	return nil
}
