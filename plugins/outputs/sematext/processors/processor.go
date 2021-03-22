package processors

import "github.com/influxdata/telegraf"

type Processor interface {
	Process(metrics []telegraf.Metric) error
}
