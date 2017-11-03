package cloudwatch

import (
	"errors"
	"time"

	"github.com/nirnanaaa/asparagus/toml"
)

const (
	// DefaultEnable is default state for metrics
	DefaultEnable = false
)

// Config represents the meta configuration.
type Config struct {
	Enabled                  bool              `toml:"enabled"`
	WebookURL                string            `toml:"webook-url"`
	MessageFormat            string            `toml:"message-format"`
	ReportingInterval        toml.Duration     `toml:"reporting-interval"`
	Namespace                string            `toml:"namespace"`
	ResetCountersAfterReport bool              `toml:"reset-counters-after-report"`
	StaticDimensions         map[string]string `toml:"static-dimensions"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() Config {
	return Config{
		Enabled:           DefaultEnable,
		ReportingInterval: toml.Duration(24 * time.Hour),
		StaticDimensions:  map[string]string{},
	}
}

func (c Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Namespace == "" {
		return errors.New("Metrics.Cloudwatch.Namespace must be specified ([metrics.cloudwatch] dir)")
	}
	return nil
}
