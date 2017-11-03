package metric

import (
	"log"
	"os"
	"time"

	"github.com/nirnanaaa/asparagus/metric/config"
	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	metrics "github.com/rcrowley/go-metrics"
)

// Service Returns a config
type Service struct {
	Config *Config
	Logger *logrus.Logger
}

func NewService(config *Config) *Service {
	return &Service{
		Config: config,
	}
}
func (s *Service) Printf(format string, v ...interface{}) {
	s.Logger.Infof(format, v)
}

// func (s *Service) Log(r metrics.Registry, freq time.Duration, l Logger) {
// LogScaled(r, freq, time.Nanosecond, l)
// }
func (s *Service) Start() error {
	// w, _ := syslog.Dial("unixgram", "/dev/log", syslog.LOG_INFO, "metrics")
	// go metrics.Syslog(metrics.DefaultRegistry, 60e9, w)

	s.Logger.Info("Starting Metrics Service")
	go metrics.Log(
		metrics.DefaultRegistry,
		21600*time.Second,
		log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	// opsQueued := prometheus.NewGaugeVec(
	// 	prometheus.GaugeOpts{
	// 		Namespace: "our_company",
	// 		Subsystem: "blob_storage",
	// 		Name:      "ops_queued",
	// 		Help:      "Number of blob storage operations waiting to be processed, partitioned by user and type.",
	// 	},
	// 	[]string{
	// 		// Which user has requested the operation?
	// 		"user",
	// 		// Of what type is the operation?
	// 		"type",
	// 	},
	// )
	// prometheus.MustRegister(opsQueued)
	// // Increase a value using compact (but order-sensitive!) WithLabelValues().
	// opsQueued.WithLabelValues("bob", "put").Add(4)
	// // Increase a value with a map using WithLabels. More verbose, but order
	// // doesn't matter anymore.
	// opsQueued.With(prometheus.Labels{"type": "delete", "user": "alice"}).Inc()
	// if err := registry.Register(prometheus.NewProcessCollector(os.Getpid(), "")); err != nil {
	// 	return err
	// }
	// if err := registry.Register(prometheus.NewProcessCollectorPIDFn(
	// 	func() (int, error) { return os.Getpid(), nil }, "foobar"),
	// ); err != nil {
	// 	return err
	// }
	//
	// _, err := registry.Gather()
	// if err != nil {
	// 	return err
	// }
	metricsRegistry := metrics.DefaultRegistry
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	}))
	registrator := &config.Config{
		Client:                cloudwatch.New(sess),
		Namespace:             "ASPARAGUS",
		Filter:                &config.NoFilter{},
		ResetCountersOnReport: true,
		ReportingInterval:     10 * time.Second,
		StaticDimensions: map[string]string{
			"Environment": os.Getenv("ENV"),
		},
	}
	go Cloudwatch(metricsRegistry, registrator)
	return nil
}

func (s *Service) Stop() error {
	return nil
}

// go metrics.Log(metrics.DefaultRegistry, 5 * time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
