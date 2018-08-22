package sql

const (
	// DefaultEnabled defines a state for default config
	DefaultEnabled = false

	// DefaultDatabaseURI defines the default contact point towards the mysql server instance
	DefaultDatabaseURI = "/tmp/asparagus.db"

	// DefaultDriver defines a default driver for the sql package
	DefaultDriver = "sqlite3"
)

// Config represents a configuration for a HTTP service.
type Config struct {
	Enabled        bool   `toml:"enabled"`
	DatabaseURI    string `toml:"uri"`
	DatabaseDriver string `toml:"driver"`
}

// NewConfig returns a new Config with default settings.
func NewConfig() Config {
	return Config{
		Enabled:        DefaultEnabled,
		DatabaseURI:    DefaultDatabaseURI,
		DatabaseDriver: DefaultDriver,
	}
}
