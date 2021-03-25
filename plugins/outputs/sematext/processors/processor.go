package processors

import "github.com/influxdata/telegraf"

// MetricProcessor is interface that should be implemented by modules which adjust Telegraf metrics to match Sematext
// format
type MetricProcessor interface {
	// Process makes adjustments to a single metric instance to be compliant with Sematext backend
	Process(metric telegraf.Metric) error
}

// BatchProcessor is used to execute actions on the level of a whole batch of metrics.
type BatchProcessor interface {
	Process(metrics []telegraf.Metric) ([]telegraf.Metric, error)
}
