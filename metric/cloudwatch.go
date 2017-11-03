package metric

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/nirnanaaa/asparagus/metric/config"
	metrics "github.com/rcrowley/go-metrics"
)

//blocks, run as go reporter.Cloudwatch(cfg)
func Cloudwatch(registry metrics.Registry, cfg *config.Config) {
	ticks := time.NewTicker(cfg.ReportingInterval)
	defer ticks.Stop()
	for {
		select {
		case <-ticks.C:
			emitMetrics(registry, cfg)
		}
	}
}

func emitMetrics(registry metrics.Registry, cfg *config.Config) {
	data := metricsData(registry, cfg)

	//20 is the max metrics per request
	for len(data) > 20 {
		put := data[0:20]
		putMetrics(cfg, put)
		data = data[20:]
	}

	if len(data) > 0 {
		putMetrics(cfg, data)
	}

}

func putMetrics(cfg *config.Config, data []*cloudwatch.MetricDatum) {
	client := cfg.Client

	req := &cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(cfg.Namespace),
		MetricData: data,
	}
	client.PutMetricData(req)
}

func metricsData(registry metrics.Registry, cfg *config.Config) []*cloudwatch.MetricDatum {

	data := []*cloudwatch.MetricDatum{}
	timestamp := aws.Time(time.Now())

	aDatum := func(name string) *cloudwatch.MetricDatum {
		return &cloudwatch.MetricDatum{
			MetricName: aws.String(name),
			Timestamp:  timestamp,
			Dimensions: dimensions(cfg),
		}
	}
	//rough port from the graphite reporter
	registry.Each(func(name string, i interface{}) {

		switch metric := i.(type) {
		case metrics.Counter:
			if metric.Count() == 0 {
				return
			}
			count := float64(metric.Count())
			if cfg.Filter.ShouldReport(name, count) {
				datum := aDatum(name)
				datum.Unit = aws.String(cloudwatch.StandardUnitCount)
				datum.Value = aws.Float64(count)
				data = append(data, datum)
			}
			if cfg.ResetCountersOnReport {
				metric.Clear()
			}
		case metrics.Gauge:
			value := float64(metric.Value())
			if cfg.Filter.ShouldReport(name, value) {
				datum := aDatum(name)
				datum.Unit = aws.String(cloudwatch.StandardUnitCount)
				datum.Value = aws.Float64(float64(value))
				data = append(data, datum)
			}
		case metrics.GaugeFloat64:
			value := float64(metric.Value())
			if cfg.Filter.ShouldReport(name, value) {
				datum := aDatum(name)
				datum.Unit = aws.String(cloudwatch.StandardUnitCount)
				datum.Value = aws.Float64(value)
				data = append(data, datum)
			}
		case metrics.Histogram:
			h := metric.Snapshot()
			value := float64(h.Count())
			if cfg.Filter.ShouldReport(name, value) {
				for _, p := range cfg.Filter.Percentiles(name) {
					pname := fmt.Sprintf("%s-perc%.3f", name, p)
					pvalue := h.Percentile(p)
					if cfg.Filter.ShouldReport(pname, pvalue) {
						datum := aDatum(pname)
						datum.Unit = aws.String(cloudwatch.StandardUnitCount)
						datum.Value = aws.Float64(pvalue)
						data = append(data, datum)
					}
				}
			}
		case metrics.Meter:
			m := metric.Snapshot()
			dataz := map[string]float64{
				fmt.Sprintf("%s.count", name):          float64(m.Count()),
				fmt.Sprintf("%s.one-minute", name):     m.Rate1(),
				fmt.Sprintf("%s.five-minute", name):    m.Rate5(),
				fmt.Sprintf("%s.fifteen-minute", name): m.Rate15(),
				fmt.Sprintf("%s.mean", name):           m.RateMean(),
			}
			for n, v := range dataz {
				if cfg.Filter.ShouldReport(n, v) {
					datum := aDatum(n)
					datum.Value = aws.Float64(v)
					data = append(data, datum)
				}
			}
		case metrics.Timer:
			t := metric.Snapshot()
			if t.Count() == 0 {
				return
			}
			dataz := map[string]float64{
				fmt.Sprintf("%s.count", name):          float64(t.Count()),
				fmt.Sprintf("%s.one-minute", name):     t.Rate1(),
				fmt.Sprintf("%s.five-minute", name):    t.Rate5(),
				fmt.Sprintf("%s.fifteen-minute", name): t.Rate15(),
				fmt.Sprintf("%s.mean", name):           t.RateMean(),
			}
			for n, v := range dataz {
				if cfg.Filter.ShouldReport(n, v) {
					datum := aDatum(n)
					datum.Value = aws.Float64(v)
					data = append(data, datum)
				}
			}
			for _, p := range cfg.Filter.Percentiles(name) {
				pname := fmt.Sprintf("%s-perc%.3f", name, p)
				pvalue := t.Percentile(p)
				if cfg.Filter.ShouldReport(pname, pvalue) {
					datum := aDatum(pname)
					datum.Value = aws.Float64(pvalue)
					data = append(data, datum)
				}
			}

		}
	})

	return data
}

func dimensions(c *config.Config) []*cloudwatch.Dimension {
	ds := []*cloudwatch.Dimension{}
	for k, v := range c.StaticDimensions {
		d := &cloudwatch.Dimension{
			Name:  aws.String(k),
			Value: aws.String(v),
		}

		ds = append(ds, d)
	}
	return ds
}
