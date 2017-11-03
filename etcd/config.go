package etcd

const (
	// DefaultServiceName is the default name to bind to.
	DefaultServiceName = "DummyService"

	// DefaultRegistryURL defines the default contact point
	DefaultRegistryURL = "http://etcd:4001"

	// DefaultLocal detects if local etcd should be used
	DefaultLocal = false

	// DefaultBindingPort describes the default port to register on
	DefaultBindingPort = 9092

	//DefaultUseDNSDiscovery defines if dns discovery should be used by default
	DefaultUseDNSDiscovery = true

	//DefaultRegistryDomain represents a domain name to resolve the SRV record from.
	DefaultRegistryDomain = "services.miha-bodytec.com"

	// DefaultCA defines a default location for the ca cert.
	DefaultCA = "/etc/ssl/certs/ca.crt"

	// DefaultUseSSL defines if ssl should be used for communicating with etcd
	DefaultUseSSL = false
)

// Config represents a configuration for a HTTP service.
type Config struct {
	ServiceName     string   `toml:"service-name"`
	RegistryURL     []string `toml:"registry-url"`
	BindingPort     int      `toml:"binding-port"`
	Local           bool     `toml:"local"`
	RegistryDomain  string   `toml:"discovery-domain"`
	UseDNSDiscovery bool     `toml:"use-dns-discovery"`
	CACert          string   `toml:"ca-cert"`
	UseSSL          bool     `toml:"use-ssl"`
}

// NewConfig returns a new Config with default settings.
func NewConfig() Config {
	return Config{
		ServiceName:     DefaultServiceName,
		RegistryURL:     []string{DefaultRegistryURL},
		Local:           DefaultLocal,
		BindingPort:     DefaultBindingPort,
		UseDNSDiscovery: DefaultUseDNSDiscovery,
		RegistryDomain:  DefaultRegistryDomain,
		CACert:          DefaultCA,
		UseSSL:          DefaultUseSSL,
	}
}
