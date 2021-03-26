package processors

import "github.com/influxdata/telegraf"

// Token is processor which injects Sematext App token into each metric
type Token struct {
	Token string
}

// Process is a method where Token processor logic is implemented
func (t *Token) Process(metric telegraf.Metric) error {
	metric.AddTag("token", t.Token)
	return nil
}

// NewToken creates a new token processor
func NewToken(token string) *Token {
	return &Token{
		Token: token,
	}
}
