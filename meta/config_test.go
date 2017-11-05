package meta_test

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/nirnanaaa/asparagus/meta"
)

func TestConfig_Parse(t *testing.T) {
	// Parse configuration.
	var c meta.Config
	if _, err := toml.Decode(`
debug = false
log-level = "debug"
logging-enabled = false
`, &c); err != nil {
		t.Fatal(err)
	}

	// Validate configuration.
	if c.Debug {
		t.Fatalf("unexpected debug: %v", c.Debug)
	} else if c.LoggingEnabled {
		t.Fatalf("unexpected logging enabled: %v", c.LoggingEnabled)
	} else if c.LogLevel != "debug" {
		t.Fatalf("unexpected log level: %v", c.LogLevel)
	}
}
