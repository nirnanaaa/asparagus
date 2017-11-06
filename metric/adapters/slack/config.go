package slack

import (
	"time"

	"github.com/nirnanaaa/asparagus/toml"
)

const (
	// DefaultEnable is default state for metrics
	DefaultEnable = false

	// DefaultLoggingEnabled is disabled by default
	DefaultLoggingEnabled = false

	// DefaultLoggingFormat defines a format for logging
	DefaultLoggingFormat = `
Cronjob *{{ .Name }}* was started at *{{ .StartDate }}* and ended at *{{ .EndDate }}*
The execution itself tool *{{ .Result.Took }}* seconds.
The resulting status code was *{{ .StatusCode }}*.
{{ if .IsError }}
The returned response was:

{{ .ResultDebug }}

{{ end }}
`
)

// Config represents the meta configuration.
type Config struct {
	Enabled              bool          `toml:"enabled"`
	DirectLoggingEnabled bool          `toml:"enable-direct-logging"`
	LogFormat            string        `toml:"logging-format"`
	WebookURL            string        `toml:"webhook-url"`
	MessageFormat        string        `toml:"message-format"`
	ReportingInterval    toml.Duration `toml:"reporting-interval"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() Config {
	return Config{
		Enabled:              DefaultEnable,
		ReportingInterval:    toml.Duration(24 * time.Hour),
		LogFormat:            DefaultLoggingFormat,
		DirectLoggingEnabled: DefaultLoggingEnabled,
	}
}

func (c Config) Validate() error {
	// if c.BasePath == "" {
	// 	return errors.New("Meta.BasePath must be specified ([meta] dir)")
	// }
	return nil
}
