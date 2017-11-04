package crontab

const (
	// DefaultEnable is default state for metrics
	DefaultEnable = true

	// DefaultSource defines a source file
	DefaultSource = "/etc/crontab"
)

// Config represents the meta configuration.
type Config struct {
	Enabled bool   `toml:"enabled"`
	Crontab string `toml:"source"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() Config {
	return Config{
		Enabled: DefaultEnable,
		Crontab: DefaultSource,
	}
}

// Validate the config
func (c Config) Validate() error {
	return nil
}
