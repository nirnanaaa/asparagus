package console

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
	ReportingInterval toml.Duration `toml:"reporting-interval"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() Config {
	return Config{
		Enabled:           DefaultEnable,
		ReportingInterval: toml.Duration(24 * time.Hour),
	}
}

// Validate the config
func (c Config) Validate() error {
	return nil
}
