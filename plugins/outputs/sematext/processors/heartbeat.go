package processors

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"sync"
	"time"
)

const (
	oneMinuteSeconds int64 = 60
)

// Heartbeat is a batch processor that injects heartbeat metric as necessary (once per minute).
type Heartbeat struct {
	lastInjectedMinute int64
	lock               sync.Mutex
}

// Process is a method where Heartbeat processor checks whether a heartbeat metric is needed and injects it if so
func (h *Heartbeat) Process(metrics []telegraf.Metric) ([]telegraf.Metric, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	now := time.Now()
	if h.heartbeatNeeded(now) {
		// a heartbeat metric will be added to the batch with "current" timestamp regardless of whether the batch
		// is a fresh one or is being resent because of an earlier failure - the only important things are that it is
		// created and that we try to send it as soon as possible
		newMetrics, err := h.addHeartbeat(metrics, now)

		if err != nil {
			return metrics, err
		}

		metrics = newMetrics
	}

	return metrics, nil
}

func (h *Heartbeat) addHeartbeat(metrics []telegraf.Metric, now time.Time) ([]telegraf.Metric, error) {
	hb, err := h.createHeartbeat(now)
	if err != nil {
		return nil, err
	}

	metrics = append(metrics, hb)
	h.lastInjectedMinute = getEpochMinute(now)

	return metrics, nil
}

func (h *Heartbeat) createHeartbeat(timestamp time.Time) (telegraf.Metric, error) {
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

func (h *Heartbeat) heartbeatNeeded(now time.Time) bool {
	nowMinute := getEpochMinute(now)
	return nowMinute > h.lastInjectedMinute
}

func getEpochMinute(time time.Time) int64 {
	return time.Unix() / oneMinuteSeconds
}
