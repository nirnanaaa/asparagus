package scheduler_test

import (
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/nirnanaaa/asparagus/scheduler"
	aToml "github.com/nirnanaaa/asparagus/toml"
)

func TestConfig_Parse(t *testing.T) {
	// Parse configuration.
	var c scheduler.Config
	if _, err := toml.Decode(`
    enabled = true
    tick-duration = "1s"
    log-task-detection = false
    num-workers = 10
    [executor-http]
      enabled = true
      debug-response = true
      log-http-status = true
      use-jwt-signing = false
      jwt-issuer = "issuer"
      jwt-expires = "10m0s"
      jwt-subject = "0"
      jwt-secret = ""
    [provider-crontab]
      enabled = true
      source = "/Users/fkasper/.asparagus/crontab"
`, &c); err != nil {
		t.Fatal(err)
	}

	// Validate configuration.
	if c.Enabled != true {
		t.Fatalf("unexpected enabled: %v", c.Enabled)
	} else if c.TickDuration != aToml.Duration(1*time.Second) {
		t.Fatalf("unexpected tick duration: %v", c.TickDuration)
	} else if c.LogTasksDetection {
		t.Fatalf("unexpected log tasks detection flag: %v", c.LogTasksDetection)
	} else if c.NumWorkers != 10 {
		t.Fatalf("expected 10 workers got %d", c.NumWorkers)
	} else if !c.HTTPExecutor.Enabled {
		t.Fatalf("unexpected http executor enabled %v", c.HTTPExecutor.Enabled)
	} else if !c.CrontabSource.Enabled {
		t.Fatalf("unexpected crontab source enabled %v", c.CrontabSource.Enabled)
	}
}
