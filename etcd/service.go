package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	uuid "github.com/satori/go.uuid"
)

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	Logger *logrus.Logger
	Config Config
	KAPI   client.KeysAPI
}

// NewService returns a new instance of Service.
func NewService(c Config) *Service {
	s := &Service{
		Logger: logrus.New(),
		Config: c,
	}
	return s
}

func formatServiceName(service string) string {
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Sprintf("/services/%s/%s", service, uuid.NewV4().String())
	}
	return fmt.Sprintf("/services/%s/%s", service, hostname)
}

// Settings are defined for an in-registry schema
type Settings struct {
	IPAddress string `json:"IPAddress"`
	Schema    string `json:"Schema"`
	Port      int    `json:"Port"`
	BaseURL   string `json:"BaseURL"`
}

// KAPISet registers the service inside etcd, every 60 seconds
func (s *Service) KAPISet() {
	opts := client.SetOptions{
		TTL: 60 * time.Second,
	}
	myIP := s.GetMYIP()
	settings := Settings{
		myIP,
		"http",
		s.Config.BindingPort,
		"",
	}
	jsonSettings, err := json.Marshal(&settings)
	if err != nil {
		s.Logger.WithError(err).Error("an error occured while marshalling the registry settings")
		return
	}
	if _, err := s.KAPI.Set(context.Background(), formatServiceName(s.Config.ServiceName), string(jsonSettings), &opts); err != nil {
		s.Logger.WithError(err).Error("an error occured while registering the services IP-Address")
		return
	}
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

// Start starts the service
func (s *Service) Start() error {
	s.Logger.Info("Starting etcd registration")

	myIP := s.GetMYIP()
	s.Logger.WithField("IPAddress", myIP).Info("Using IP Address")

	connectIP := []string{fmt.Sprintf("http://%s:4001", myIP)}
	if s.Config.Local {
		connectIP = s.Config.RegistryURL
	}
	if s.Config.UseDNSDiscovery {
		discoverer := client.NewSRVDiscover()
		endpoints, err := discoverer.Discover(s.Config.RegistryDomain)
		if err != nil {
			return err
		}
		connectIP = endpoints
	}
	s.Logger.WithField("Endpoints", connectIP).Debug("Discovery Endpoints")
	tlsConfig := &tls.Config{}
	if s.Config.UseSSL {
		// Load client cert
		caCert, err := ioutil.ReadFile(s.Config.CACert)
		if err != nil {
			return err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
		tlsConfig.BuildNameToCertificate()
	}
	t := http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     tlsConfig,
	}
	cfg := client.Config{
		Endpoints: connectIP,
		Transport: &t,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return err
	}
	s.KAPI = client.NewKeysAPI(c)
	go func() {
		for {
			s.KAPISet()
			time.Sleep(50 * time.Second)
		}
	}()
	return nil
}

// GetMYIP returns the servers IP Address
func (s *Service) GetMYIP() string {
	client := http.Client{
		Timeout: 500 * time.Millisecond,
	}
	internalIP, err := externalIP()
	if err != nil {
		s.Logger.WithError(err).Error("Error while getting local IP Address. Trying next method.")
		return os.Getenv("ASPARAGUS_ETCD_SERVICE_NAME")
	}
	resp, err := client.Get("http://169.254.169.254/latest/meta-data/local-ipv4")
	if err != nil {
		return internalIP
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return internalIP
	}
	return string(body)
}

// Stop closes the underlying listener.
func (s *Service) Stop() error {
	return nil
}
