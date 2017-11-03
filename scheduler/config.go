package scheduler

const (
	//DefaultTickDuration is default state for metrics
	DefaultTickDuration = 1000

	//DefaultAPIFallbackDomain defines a default fallback
	DefaultAPIFallbackDomain = "example.com"
)

// Config represents the meta configuration.
type Config struct {
	TickDuration      int64  `toml:"tick-duration-ms"`
	APIFallbackDomain string `toml:"api-fallback-domain"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		TickDuration:      DefaultTickDuration,
		APIFallbackDomain: DefaultAPIFallbackDomain,
	}
}

func (c *Config) Validate() error {
	return nil
}
