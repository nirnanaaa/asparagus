package local_test

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/nirnanaaa/asparagus/scheduler/provider/local"
)

func TestConfig_Parse(t *testing.T) {
	// Parse configuration.
	var c local.Config
	if _, err := toml.Decode(`
    enabled = true
		enable-output = false
`, &c); err != nil {
		t.Fatal(err)
	}

	// Validate configuration.
	if c.Enabled != true {
		t.Fatalf("unexpected enabled: %v", c.Enabled)
	} else if c.EnabledOutput {
		t.Fatalf("unexpected enabled output: %v", c.EnabledOutput)
	}
}
