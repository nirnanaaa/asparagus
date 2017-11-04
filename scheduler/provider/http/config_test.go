package http_test

import (
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/nirnanaaa/asparagus/scheduler/provider/http"
	aToml "github.com/nirnanaaa/asparagus/toml"
)

func TestConfig_Parse(t *testing.T) {
	// Parse configuration.
	var c http.Config
	if _, err := toml.Decode(`
    enabled = true
    debug-response = true
    log-http-status = true
    use-jwt-signing = false
    jwt-issuer = "issuer"
    jwt-expires = "10m0s"
    jwt-subject = "0"
    jwt-secret = ""
`, &c); err != nil {
		t.Fatal(err)
	}

	// Validate configuration.
	if c.Enabled != true {
		t.Fatalf("unexpected enabled: %v", c.Enabled)
	} else if c.DebugResponse != true {
		t.Fatalf("unexpected debug: %v", c.DebugResponse)
	} else if c.LogHTTPStatus != true {
		t.Fatalf("unexpected http status logging: %v", c.LogHTTPStatus)
	} else if c.SignJWT {
		t.Fatalf("unexpected sign jwt: %v", c.SignJWT)
	} else if c.JWTSubject != "0" {
		t.Fatalf("unexpected jwt subject: %v", c.JWTSubject)
	} else if c.Issuer != "issuer" {
		t.Fatalf("unexpected issuer: %v", c.Issuer)
	} else if c.JWTSecret != "" {
		t.Fatalf("unexpected secret: %v", c.JWTSecret)
	} else if c.JWTExpires != aToml.Duration(10*time.Minute) {
		t.Fatalf("unexpected expire time on jwt: %v", c.JWTExpires)
	}
}
