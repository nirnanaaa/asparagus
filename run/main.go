package run

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nirnanaaa/asparagus/metric"
	"github.com/nirnanaaa/asparagus/scheduler"
)

// Service interface definition
type Service interface {
	Start() error
	Stop() error
}

// Command represents the command executed by "asparagus run".
type Command struct {
	Version   string
	Branch    string
	Commit    string
	BuildTime string

	closing chan struct{}
	err     chan error
	Closed  chan struct{}

	Stdin    io.Reader
	Stdout   io.Writer
	Stderr   io.Writer
	Services []Service
}

// NewCommand return a new instance of Command.
func NewCommand() *Command {
	return &Command{
		closing: make(chan struct{}),
		err:     make(chan error),
		Closed:  make(chan struct{}),
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}
}

const logo = `
    ___
   /   |  _________  ____ __________ _____ ___  _______
  / /| | / ___/ __ \/ __  / ___/ __ \/ __ \/ / / / ___/
 / ___ |(__  ) /_/ / /_/ / /  / /_/ / /_/ / /_/ (__  )
/_/  |_/____/ .___/\__,_/_/   \__,_/\__, /\__,_/____/
           /_/                     /____/
`

// Run parses the config from args and runs the server.
func (cmd *Command) Run(args ...string) error {
	options, err := cmd.ParseFlags(args...)
	if err != nil {
		return err
	}
	logger := log.New()

	// Print sweet asparagus logo.
	fmt.Print(logo)

	// Set parallelism.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Mark start-up in log.
	logger.Infof("Asparagus starting, version %s, branch %s, commit %s",
		cmd.Version, cmd.Branch, cmd.Commit)
	logger.Infof("Go version %s, GOMAXPROCS set to %d", runtime.Version(), runtime.GOMAXPROCS(0))
	runtime.SetBlockProfileRate(int(1 * time.Second))
	config, err := cmd.ParseConfig(options.GetConfigPath())
	if err != nil {
		return fmt.Errorf("parse config: %s", err)
	}
	lvl, err := log.ParseLevel(config.Meta.LogLevel)
	if err != nil {
		return fmt.Errorf("parse log level: %s", err)
	}
	log.SetLevel(lvl) //.config.Meta.LogLevel)

	// Apply any environment variables on top of the parsed config
	if err := config.ApplyEnvOverrides(); err != nil {
		return fmt.Errorf("apply env config: %v", err)
	}

	// Validate the configuration.
	if err := config.Validate(); err != nil {
		return fmt.Errorf("%s. To generate a valid configuration file run `asparagus config > asparagus.generated.conf`", err)
	}
	flag.Parse()

	metricx := metric.NewService(config.Metrics)
	metricx.Logger = logger
	cmd.Services = append(cmd.Services, metricx)
	scheduler := scheduler.NewService(config.Scheduler)
	scheduler.Logger = logger
	cmd.Services = append(cmd.Services, scheduler)
	return nil
}

// Open opens all services
func (cmd *Command) Open() error {
	for _, service := range cmd.Services {
		if err := service.Start(); err != nil {
			log.Errorf(err.Error())
			return fmt.Errorf("open service: %s", err)
		}
	}
	return nil
}

// Close shuts down the server.
func (cmd *Command) Close() error {
	// Close services to allow any inflight requests to complete
	// and prevent new requests from being accepted.
	for _, service := range cmd.Services {
		service.Stop()
	}
	defer close(cmd.Closed)
	close(cmd.closing)
	return nil
}

// Options represents the command line options that can be parsed.
type Options struct {
	ConfigPath string
	Action     string
}

var usage = `Runs the Asparagus daemon.
Usage: asparagus <cmd> run [flags]
    -config <path>
            Set the path to the configuration file.
            This defaults to the environment variable ASPARAGUS_CONFIG_PATH,
            ~/.asparagus/asparagus.conf, or /etc/asparagus/asparagus.conf if a file
            is present at any of these locations.
            Disable the automatic loading of a configuration file using
            the null device (such as /dev/null).
`

// ParseFlags parses the command line flags from args and returns an options set.
func (cmd *Command) ParseFlags(args ...string) (Options, error) {
	var options Options
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&options.ConfigPath, "config", "", "")
	fs.StringVar(&options.Action, "s", "", "")
	_ = fs.String("hostname", "", "")
	fs.Usage = func() {
		fmt.Fprintln(cmd.Stderr, usage)
	}
	if err := fs.Parse(args); err != nil {
		return Options{}, err
	}
	return options, nil
}

// ParseConfig parses the config at path.
// Returns a demo configuration if path is blank.
func (cmd *Command) ParseConfig(path string) (*Config, error) {
	// Use demo configuration if no config path is specified.
	if path == "" {
		log.Println("no configuration provided, using default settings")
		return NewDemoConfig()
	}

	log.Printf("Using configuration at: %s\n", path)

	config := NewConfig()
	if err := config.FromTomlFile(path); err != nil {
		return nil, err
	}

	return config, nil
}

// GetConfigPath returns the config path from the options.
// It will return a path by searching in this order:
//   1. The CLI option in ConfigPath
//   2. The environment variable ASPARAGUS_CONFIG_PATH
//   3. The first asparagus.conf file on the path:
//        - ~/.asparagus
//        - /etc/asparagus
func (opt *Options) GetConfigPath() string {
	if opt.ConfigPath != "" {
		if opt.ConfigPath == os.DevNull {
			return ""
		}
		return opt.ConfigPath
	} else if envVar := os.Getenv("ASPARAGUS_CONFIG_PATH"); envVar != "" {
		return envVar
	}

	for _, path := range []string{
		os.ExpandEnv("${HOME}/.asparagus/asparagus.conf"),
		"/etc/asparagus/asparagus.conf",
	} {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}
