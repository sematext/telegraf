package processors

import "github.com/influxdata/telegraf"

type Token struct {
	Token string
}

func (t *Token) Process(metrics []telegraf.Metric) error {
	for _, m := range metrics {
		m.AddTag("token", t.Token)
	}
	return nil
}
