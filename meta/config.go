package meta

const (
	// DefaultLoggingEnabled determines if log messages are printed for the meta service
	DefaultLoggingEnabled = true

	// DefaultLogFile defines the destination for logs
	DefaultLogFile = "/var/log/queue.log"
)

// Config represents the meta configuration.
type Config struct {
	Debug          bool   `toml:"debug"`
	LogFile        string `toml:"log-file"`
	LoggingEnabled bool   `toml:"logging-enabled"`
	Version        string
}

// NewConfig builds a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		LoggingEnabled: DefaultLoggingEnabled,
		LogFile:        DefaultLogFile,
		Debug:          false,
	}
}

func (c *Config) Validate() error {
	// if c.BasePath == "" {
	// 	return errors.New("Meta.BasePath must be specified ([meta] dir)")
	// }
	return nil
}
