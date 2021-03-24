package processors

import "github.com/influxdata/telegraf"

// Token is processor which injects Sematext App token into each metric
type Token struct {
	Token string
}

// Process is a method where Token processor logic is implemented
func (t *Token) Process(metrics []telegraf.Metric) error {
	for _, m := range metrics {
		m.AddTag("token", t.Token)
	}
	return nil
}
