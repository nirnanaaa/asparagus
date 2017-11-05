package local

const (
	// DefaultEnable is default state for local execution
	DefaultEnable = false

	//DefaultEnableOutput enabled stdout logging
	DefaultEnableOutput = false
)

// Config represents the meta configuration.
type Config struct {
	Enabled       bool `toml:"enabled"`
	EnabledOutput bool `toml:"enable-output"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() Config {
	return Config{
		Enabled:       DefaultEnable,
		EnabledOutput: DefaultEnableOutput,
	}
}

// Validate the config
func (c Config) Validate() error {
	return nil
}
