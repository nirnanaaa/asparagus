package slack

const (
	// DefaultEnable is default state for metrics
	DefaultEnable = false
)

// Config represents the meta configuration.
type Config struct {
	Enabled       bool   `toml:"enabled"`
	WebookURL     string `toml:"webook-url"`
	MessageFormat string `toml:"message-format"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		Enabled: DefaultEnable,
	}
}

func (c *Config) Validate() error {
	// if c.BasePath == "" {
	// 	return errors.New("Meta.BasePath must be specified ([meta] dir)")
	// }
	return nil
}
