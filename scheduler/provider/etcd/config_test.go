package etcd_test

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/nirnanaaa/asparagus/scheduler/provider/etcd"
)

func TestConfig_Parse(t *testing.T) {
	// Parse configuration.
	c := etcd.NewConfig()
	if _, err := toml.Decode(`
enabled = false
registry-url = ["http://127.0.0.1:4001"]
discovery-enabled = false
discovery-domain = "discover.example.com"
connect-ca-cert = "/etc/ssl/certs/ca.crt"
source-folder = "/cron/jobs"
json-single-depth = false
`, &c); err != nil {
		t.Fatal(err)
	}

	// Validate configuration.
	if c.Enabled {
		t.Fatalf("unexpected enabled: %v", c.Enabled)
	} else if c.UseJSONSingleDepth {
		t.Fatalf("unexpected json depth: %v", c.UseJSONSingleDepth)
	} else if c.UseDNSDiscovery {
		t.Fatalf("unexpected dns discovery flag: %v", c.UseDNSDiscovery)
	} else if c.DiscoveryDomain != "discover.example.com" {
		t.Fatalf("unexpected discovery domain: %v", c.DiscoveryDomain)
	} else if len(c.RegistryURL) != 1 {
		t.Fatalf("unexpected registry urls: %v", c.RegistryURL)
	} else if c.CACert != "/etc/ssl/certs/ca.crt" {
		t.Fatalf("unexpected ssl cert root: %v", c.CACert)
	} else if c.SourceFolder != "/cron/jobs" {
		t.Fatalf("unexpected source folder: %v", c.SourceFolder)
	}
}
