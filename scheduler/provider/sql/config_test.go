package sql_test

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/nirnanaaa/asparagus/scheduler/provider/sql"
)

func TestConfig_Parse(t *testing.T) {
	// Parse configuration.
	c := sql.NewConfig()
	if _, err := toml.Decode(`
enabled = false
uri = "/tmp/test.db"
driver = "test"
`, &c); err != nil {
		t.Fatal(err)
	}

	// Validate configuration.
	if c.Enabled {
		t.Fatalf("unexpected enabled: %v", c.Enabled)
	} else if c.DatabaseURI != "/tmp/test.db" {
		t.Fatalf("unexpected database uri: %v", c.DatabaseURI)
	} else if c.DatabaseDriver != "test" {
		t.Fatalf("unexpected database driver: %v", c.DatabaseDriver)
	}
}
