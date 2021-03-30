package processors

import (
	"fmt"
	"github.com/influxdata/telegraf"
	"sync"
)

// MetricMetainfo contains metainfo about a single metric
type MetricMetainfo struct {
	token       string
	name        string
	namespace   string
	semType     SematextMetricType
	numericType NumericType
	label       string
	description string
	host        string
}

// SematextMetricType is an enumeration of metric types expected by Sematext backend
type SematextMetricType int

// Possible values for the ValueType enum.
const (
	_ SematextMetricType = iota
	Counter
	Gauge
)

// NumericType represents metric's data type
type NumericType int

const (
	UnsupportedNumericType NumericType = iota
	Long
	Double
	Bool
)

// Metainfo is a processor that extracts metainfo from telegraf metrics and sends it to Sematext backend
type Metainfo struct {
	log         telegraf.Logger
	token       string
	sentMetrics map[string]*MetricMetainfo
	lock        sync.Mutex
}

// NewMetainfo creates a new Metainfo processor
func NewMetainfo(log telegraf.Logger, token string) *Metainfo {
	sentMetricsMap := make(map[string]*MetricMetainfo)
	return &Metainfo{
		log:         log,
		token:       token,
		sentMetrics: sentMetricsMap,
	}
}

// Process contains core logic of Metainfo processor
func (m *Metainfo) Process(metrics []telegraf.Metric) ([]telegraf.Metric, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	newMetrics := make([]*MetricMetainfo, 0)

	for _, metric := range metrics {
		for _, field := range metric.FieldList() {
			mInfo := processMetric(m.token, &metric, field, &m.sentMetrics)
			if mInfo != nil {
				newMetrics = append(newMetrics, mInfo)
			}
		}
	}

	if len(newMetrics) > 0 {
		// TODO send the metainfo
		// TODO add newMetrics to sentMetrics (maybe right away, don't wait for goroutine to finish because of locks)
	}

	return metrics, nil
}

func processMetric(token string, metric *telegraf.Metric, field *telegraf.Field,
	sentMetrics *map[string]*MetricMetainfo) *MetricMetainfo {
	host, set := (*metric).GetTag(telegrafHostTag)
	// skip if no host tag
	if set {
		key := buildMetricKey(host, (*metric).Name(), field.Key)

		_, set := (*sentMetrics)[key]
		if !set {
			return buildMetainfo(token, host, metric, field)
		}
	}

	return nil
}

func buildMetainfo(token string, host string, metric *telegraf.Metric, field *telegraf.Field) *MetricMetainfo {
	semType := getSematextMetricType((*metric).Type())
	numericType := getSematextNumericType(field)

	if numericType == UnsupportedNumericType {
		return nil
	}

	label := fmt.Sprintf("%s.%s", (*metric).Name(), field.Key)

	return &MetricMetainfo{
		token:       token,
		name:        field.Key,
		namespace:   (*metric).Name(),
		semType:     semType,
		numericType: numericType,
		label:       label,
		description: "",
		host:        host,
	}
}

// Close clears the resources used by Metainfo processor
func (m *Metainfo) Close() {
	// TODO close http client
}

func getSematextMetricType(metricType telegraf.ValueType) SematextMetricType {
	var semType SematextMetricType
	switch metricType {
	case telegraf.Counter:
		semType = Counter
	default:
		semType = Gauge
	}

	return semType
}

func getSematextNumericType(field *telegraf.Field) NumericType {
	var numType NumericType
	switch field.Value.(type) {
	case float64:
		numType = Double
	case uint64:
		numType = Long
	case int64:
		numType = Long
	case bool:
		numType = Bool
	default:
		numType = UnsupportedNumericType
	}

	return numType
}

func buildMetricKey(host string, namespace string, name string) string {
	return fmt.Sprintf("%s-%s.%s", host, namespace, name)
}
