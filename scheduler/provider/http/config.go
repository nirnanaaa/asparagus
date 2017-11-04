package http

import (
	"time"

	"github.com/nirnanaaa/asparagus/toml"
)

const (
	// DefaultEnable is default state for metrics
	DefaultEnable = true

	// DefaultDebugResponse let's you debug your HTTP Responses
	DefaultDebugResponse = false

	// DefaultLogStatusCodes let's you debug your HTTP Responses
	DefaultLogStatusCodes = false

	// DefaultSignJWT sign requests using jwt
	DefaultSignJWT = true

	// DefaultIssuer for jwt
	DefaultIssuer = "issuer"

	// DefaultJWTSubject for jwt
	DefaultJWTSubject = "0"
)

// Config represents the meta configuration.
type Config struct {
	Enabled       bool          `toml:"enabled"`
	DebugResponse bool          `toml:"debug-response"`
	LogHTTPStatus bool          `toml:"log-http-status"`
	SignJWT       bool          `toml:"use-jwt-signing"`
	Issuer        string        `toml:"jwt-issuer"`
	JWTExpires    toml.Duration `toml:"jwt-expires"`
	JWTSubject    string        `toml:"jwt-subject"`
	JWTSecret     string        `toml:"jwt-secret"`
}

// NewConfig builds a new configuration with default values.
func NewConfig() Config {
	return Config{
		Enabled:       DefaultEnable,
		DebugResponse: DefaultDebugResponse,
		LogHTTPStatus: DefaultLogStatusCodes,
		SignJWT:       DefaultSignJWT,
		Issuer:        DefaultIssuer,
		JWTExpires:    toml.Duration(10 * time.Minute),
		JWTSubject:    DefaultJWTSubject,
	}
}

// Validate the config
func (c Config) Validate() error {
	return nil
}
