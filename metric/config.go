package metric

const (
	//DefaultMetricsEnable is default state for metrics
	DefaultMetricsEnable = false
)

// Config represents the meta configuration.
type Config struct {
	Enabled bool `toml:"enable-metrics"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		Enabled: DefaultMetricsEnable,
	}
}

func (c *Config) Validate() error {
	// if c.BasePath == "" {
	// 	return errors.New("Meta.BasePath must be specified ([meta] dir)")
	// }
	return nil
}
