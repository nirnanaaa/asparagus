package etcd

const (
	// DefaultEnabled defines a state for default config
	DefaultEnabled = false

	// DefaultRegistryURL defines the default contact point
	DefaultRegistryURL = "http://127.0.0.1:4001"

	//DefaultUseDNSDiscovery defines if dns discovery should be used by default
	DefaultUseDNSDiscovery = false

	//DefaultDiscoveryDomain represents a domain name to resolve the SRV record from.
	DefaultDiscoveryDomain = "example.com"

	// DefaultCA defines a default location for the ca cert.
	DefaultCA = ""

	// DefaultSourceFolder defines a safe default for job sources
	DefaultSourceFolder = "/cron/Jobs"

	// DefaultJSONSingleDepth let's you use a flat directory structure inside etcd
	// instead of the file based one (could be a bit faster and less heavy on your network)
	DefaultJSONSingleDepth = true
)

// Config represents a configuration for a HTTP service.
type Config struct {
	Enabled            bool     `toml:"enabled"`
	RegistryURL        []string `toml:"registry-url"`
	UseDNSDiscovery    bool     `toml:"discovery-enabled"`
	DiscoveryDomain    string   `toml:"discovery-domain"`
	CACert             string   `toml:"connect-ca-cert"`
	SourceFolder       string   `toml:"source-folder"`
	UseJSONSingleDepth bool     `toml:"json-single-depth"`
}

// NewConfig returns a new Config with default settings.
func NewConfig() Config {
	return Config{
		Enabled:            DefaultEnabled,
		RegistryURL:        []string{DefaultRegistryURL},
		UseDNSDiscovery:    DefaultUseDNSDiscovery,
		DiscoveryDomain:    DefaultDiscoveryDomain,
		CACert:             DefaultCA,
		SourceFolder:       DefaultSourceFolder,
		UseJSONSingleDepth: DefaultJSONSingleDepth,
	}
}
