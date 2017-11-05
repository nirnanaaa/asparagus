package meta

const (
	// DefaultLoggingEnabled determines if log messages are printed for the meta service
	DefaultLoggingEnabled = true
)

// Config represents the meta configuration.
type Config struct {
	Debug          bool   `toml:"debug"`
	LoggingEnabled bool   `toml:"logging-enabled"`
	LogLevel       string `toml:"log-level"`
	Version        string `toml:"-"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		LoggingEnabled: DefaultLoggingEnabled,
		Debug:          false,
		LogLevel:       "warning",
	}
}

// Validate the config
func (c *Config) Validate() error {
	// if c.BasePath == "" {
	// 	return errors.New("Meta.BasePath must be specified ([meta] dir)")
	// }
	return nil
}
