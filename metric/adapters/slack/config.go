package slack

import (
	"time"

	"github.com/nirnanaaa/asparagus/toml"
)

const (
	// DefaultEnable is default state for metrics
	DefaultEnable = false
)

// Config represents the meta configuration.
type Config struct {
	Enabled           bool          `toml:"enabled"`
	WebookURL         string        `toml:"webhook-url"`
	MessageFormat     string        `toml:"message-format"`
	ReportingInterval toml.Duration `toml:"reporting-interval"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() Config {
	return Config{
		Enabled:           DefaultEnable,
		ReportingInterval: toml.Duration(24 * time.Hour),
	}
}

func (c Config) Validate() error {
	// if c.BasePath == "" {
	// 	return errors.New("Meta.BasePath must be specified ([meta] dir)")
	// }
	return nil
}
