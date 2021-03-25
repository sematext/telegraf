package processors

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"time"
)

const (
	oneMinuteSeconds int64 = 60
)

// Heartbeat is a batch processor that injects heartbeat metric as necessary (once per minute)
type Heartbeat struct {
	lastInjectedMinute int64
}

// Process is a method where Heartbeat processor logic is implemented
func (t *Heartbeat) Process(metrics []telegraf.Metric) ([]telegraf.Metric, error) {
	now := time.Now()
	if t.heartbeatNeeded(now) {
		newMetrics, err := t.addHeartbeat(metrics, now)

		if err != nil {
			return metrics, err
		}

		metrics = newMetrics
	}

	return metrics, nil
}

func (t *Heartbeat) addHeartbeat(metrics []telegraf.Metric, now time.Time) ([]telegraf.Metric, error) {
	hb, err := t.createHeartbeat(now)
	if err != nil {
		return nil, err
	}

	metrics = append(metrics, hb)
	t.lastInjectedMinute = getEpochMinute(now)

	return metrics, nil
}

func (t *Heartbeat) createHeartbeat(timestamp time.Time) (telegraf.Metric, error) {
	// no need to inject any Sematext specific tags since MetricProcessors will be run afterwards and will take care
	// of such things
	hb, err := metric.New("heartbeat",
		make(map[string]string),
		map[string]interface{}{"alive": int64(1)},
		timestamp, telegraf.Gauge)

	if err != nil {
		return nil, err
	}

	return hb, nil
}

func (t *Heartbeat) heartbeatNeeded(now time.Time) bool {
	nowMinute := getEpochMinute(now)
	return nowMinute > t.lastInjectedMinute
}

func getEpochMinute(time time.Time) int64 {
	return time.Unix() / oneMinuteSeconds
}
