package processors

import "github.com/influxdata/telegraf"

// Processor is interface that should be implemented by modules which adjust Telegraf metrics to match Sematext format
type Processor interface {
	Process(metric telegraf.Metric) error
}
