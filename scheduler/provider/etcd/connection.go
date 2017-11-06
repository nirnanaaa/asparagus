package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/nirnanaaa/asparagus/scheduler/provider"
)

// Connection defines a wrapper for the actual etcd connection
type Connection struct {
	Config Config
	KAPI   client.KeysAPI
}

// NewConnection creates a new etcd connection
func NewConnection(config Config) *Connection {
	return &Connection{
		Config: config,
	}
}

// Connect connects the etcd target
func (c *Connection) Connect() error {
	tlsConfig := &tls.Config{}
	if c.Config.CACert != "" {
		caCert, err := ioutil.ReadFile(c.Config.CACert)
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
	connect, err := c.getETCDConnectionURL()
	if err != nil {
		return err
	}
	cfg := client.Config{
		Endpoints: connect,
		Transport: &t,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	conn, err := client.New(cfg)
	if err != nil {
		return err
	}
	c.KAPI = client.NewKeysAPI(conn)
	return nil
}

// WriteTask persists a task inside etcd
func (c *Connection) WriteTask(t *provider.Task) error {
	jsonSettings, err := json.Marshal(t)
	if err != nil {
		return err
	}
	opts := client.SetOptions{}
	if _, err := c.KAPI.Set(context.Background(), c.pathForKey(t.Key), string(jsonSettings), &opts); err != nil {
		return err
	}
	return nil
}

func (c *Connection) assignKeyToTask(path string, t *provider.Task) error {
	t.Key = filepath.Base(path)
	return nil
}

// ReadTask reads a task from etcd
func (c *Connection) ReadTask(key string, t *provider.Task) error {
	return c.readTaskFullDir(c.pathForKey(key), t)
}

func (c *Connection) readTaskFullDir(path string, t *provider.Task) error {
	opts := client.GetOptions{}
	node, err := c.KAPI.Get(context.Background(), path, &opts)
	if err != nil {
		return err
	}
	val := node.Node.Value
	if err := c.parseTask(val, t); err != nil {
		return err
	}
	if err := c.assignKeyToTask(path, t); err != nil {
		return err
	}
	return nil
}

func (c *Connection) parseTask(input string, t *provider.Task) error {
	if err := json.Unmarshal([]byte(input), t); err != nil {
		return err
	}
	return nil
}

// ReadAll reads all cronjobs inside a directory
func (c *Connection) ReadAll(tasks *[]provider.Task) error {
	opts := client.GetOptions{Recursive: true}
	nodes, err := c.KAPI.Get(context.Background(), c.Config.SourceFolder, &opts)
	if err != nil {
		return err
	}
	var intermediate []provider.Task
	for _, node := range nodes.Node.Nodes {
		var t provider.Task
		if err := c.parseTask(node.Value, &t); err != nil {
			return err
		}
		if err := c.assignKeyToTask(node.Key, &t); err != nil {
			return err
		}
		intermediate = append(intermediate, t)
	}
	*tasks = intermediate
	return nil
}

// WatchTasks looks for new keys inside etcd
func (c *Connection) WatchTasks(cb func(t *provider.Task)) {
	watcherOpts := client.WatcherOptions{AfterIndex: 0, Recursive: true}
	w := c.KAPI.Watcher(c.Config.SourceFolder, &watcherOpts)
	for {
		item, err := w.Next(context.Background())
		if err != nil {
			// TODO: improve here
			continue
		}
		var task provider.Task
		if err := c.readTaskFullDir(item.Node.Key, &task); err != nil {
			// TODO: improve here
			continue
		}
		cb(&task)
	}
}

func (c *Connection) pathForKey(key string) string {
	return filepath.Join(c.Config.SourceFolder, key)
}

func (c *Connection) getETCDConnectionURL() ([]string, error) {
	domains := c.Config.RegistryURL
	if c.Config.UseDNSDiscovery {
		discoverer := client.NewSRVDiscover()
		endpoints, err := discoverer.Discover(c.Config.DiscoveryDomain)
		if err != nil {
			return domains, err
		}
		return endpoints, nil
	}
	return domains, nil
}
