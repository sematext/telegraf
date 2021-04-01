package serializer

import (
	"fmt"
	"github.com/influxdata/telegraf/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
)

func TestWrite(t *testing.T) {
	serializer := NewLinePerMetricSerializer(testutil.Logger{})

	now := time.Now()

	m, _ := metric.New(
		"os",
		map[string]string{"os.host": "hostname"},
		map[string]interface{}{"disk.size": uint64(777)},
		now)

	metrics := []telegraf.Metric{m}

	assert.Equal(t,
		fmt.Sprintf("os,os.host=hostname disk.size=777i %d\n", now.UnixNano()),
		string(serializer.Write(metrics)))

	m, _ = metric.New(
		"system",
		map[string]string{"os.host": "hostname", "token": "token"},
		map[string]interface{}{"uptime_format": "18 days, 22:37"},
		now)

	metrics = []telegraf.Metric{m}

	assert.Equal(t,
		fmt.Sprintf("system,os.host=hostname,token=token uptime_format=\"18 days, 22:37\" %d\n", now.UnixNano()),
		string(serializer.Write(metrics)))
}

func TestWriteNoTags(t *testing.T) {
	serializer := NewLinePerMetricSerializer(testutil.Logger{})

	now := time.Now()

	m, _ := metric.New(
		"os",
		map[string]string{},
		map[string]interface{}{"disk.size": uint64(777)},
		now)

	metrics := []telegraf.Metric{m}

	assert.Equal(t,
		fmt.Sprintf("os disk.size=777i %d\n", now.UnixNano()),
		string(serializer.Write(metrics)))
}

func TestWriteNoMetrics(t *testing.T) {
	serializer := NewLinePerMetricSerializer(testutil.Logger{})

	now := time.Now()

	m, _ := metric.New(
		"os",
		map[string]string{"os.host": "hostname"},
		map[string]interface{}{},
		now)

	metrics := []telegraf.Metric{m}

	assert.Equal(t, "", string(serializer.Write(metrics)))
}

func TestWriteMultipleTagsAndMetrics(t *testing.T) {
	serializer := NewLinePerMetricSerializer(testutil.Logger{})

	now := time.Now()

	m, _ := metric.New(
		"os",
		map[string]string{"os.host": "hostname", "os.disk": "sda1"},
		map[string]interface{}{"disk.used": float64(12.34), "disk.free": int64(55), "disk.size": uint64(777)},
		now)

	metrics := []telegraf.Metric{m}

	assert.Equal(t,
		fmt.Sprintf("os,os.disk=sda1,os.host=hostname disk.free=55i,disk.size=777i,disk.used=12.34 %d\n", now.UnixNano()),
		string(serializer.Write(metrics)))
}
