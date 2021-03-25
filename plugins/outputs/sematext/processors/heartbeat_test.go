package processors

import (
	"github.com/influxdata/telegraf"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateHeartbeat(t *testing.T) {
	h := &Heartbeat{}

	now := time.Now()
	metric, err := h.createHeartbeat(now)

	assert.Nil(t, err)
	assert.Equal(t, "heartbeat", metric.Name())
	assert.Equal(t, 1, len(metric.Fields()))
	val, set := metric.GetField("alive")
	assert.Equal(t, true, set)
	assert.Equal(t, int64(1), val)
	assert.Equal(t, 0, len(metric.Tags()))
}

func TestHeartbeatNeeded(t *testing.T) {
	now := time.Now()
	currentMinute := getEpochMinute(now)
	h := &Heartbeat{}
	assert.Equal(t, true, h.heartbeatNeeded(now))

	h.lastInjectedMinute = currentMinute
	assert.Equal(t, false, h.heartbeatNeeded(now))

	h.lastInjectedMinute = currentMinute - 1
	assert.Equal(t, true, h.heartbeatNeeded(now))

	h.lastInjectedMinute = currentMinute + 1
	assert.Equal(t, false, h.heartbeatNeeded(now))
}

func TestAddHeartbeat(t *testing.T) {
	now := time.Now()
	currentMinute := getEpochMinute(now)
	h := &Heartbeat{
		lastInjectedMinute: currentMinute - 1,
	}

	metrics := make([]telegraf.Metric, 0, 1)
	var err error
	metrics, err = h.addHeartbeat(metrics, now)

	assert.Nil(t, err)
	assert.Equal(t, currentMinute, h.lastInjectedMinute)
	assert.Equal(t, 1, len(metrics))
}

func TestProcess(t *testing.T) {
	h := &Heartbeat{}
	metrics := make([]telegraf.Metric, 0, 1)

	assert.Equal(t, int64(0), h.lastInjectedMinute)

	var err error
	metrics, err = h.Process(metrics)

	assert.Nil(t, err)
	assert.NotEqual(t, int64(0), h.lastInjectedMinute)
	assert.Equal(t, 1, len(metrics))
}
