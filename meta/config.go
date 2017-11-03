package meta

import "github.com/Sirupsen/logrus"

const (
	// DefaultLoggingEnabled determines if log messages are printed for the meta service
	DefaultLoggingEnabled = true
)

// Config represents the meta configuration.
type Config struct {
	Debug          bool         `toml:"debug"`
	LoggingEnabled bool         `toml:"logging-enabled"`
	LogLevel       logrus.Level `toml:"log-level"`
	Version        string
}

// NewConfig builds a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		LoggingEnabled: DefaultLoggingEnabled,
		Debug:          false,
		LogLevel:       logrus.WarnLevel,
	}
}

func (c *Config) Validate() error {
	// if c.BasePath == "" {
	// 	return errors.New("Meta.BasePath must be specified ([meta] dir)")
	// }
	return nil
}
