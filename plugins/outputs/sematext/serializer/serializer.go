package serializer

import (
	"bytes"
	"fmt"
	"github.com/influxdata/telegraf"
	"math"
	"sort"
	"strconv"
	"strings"
)

// MetricSerializer is an interface implemented by different metric serialization implementations.
type MetricSerializer interface {
	// Write serializes metrics from metrics parameter and returns the result in []byte
	Write(metrics []telegraf.Metric) []byte
}

// LinePerMetricSerializer provides simple implementation which writes each metric into a new line. The logic is simpler
// and lighter, but the resulting output will be bigger.
type LinePerMetricSerializer struct {
	log telegraf.Logger
}

// NewLinePerMetricSerializer creates an instance of NewLinePerMetricSerializer
func NewLinePerMetricSerializer(log telegraf.Logger) *LinePerMetricSerializer {
	return &LinePerMetricSerializer{
		log: log,
	}
}

// NewMetricSerializer creates and instance serializer which should be used to produce Sematext metrics format
func NewMetricSerializer(log telegraf.Logger) MetricSerializer {
	return NewLinePerMetricSerializer(log)
}

// Write serializes input metrics array according to Sematext variant of influx line protocol. The output is returned
// as []byte which can be empty if there were no metrics or metrics couldn't be serialized.
func (s *LinePerMetricSerializer) Write(metrics []telegraf.Metric) []byte {
	var output bytes.Buffer

	// sematext format is based on influx line protocol: namespace,tags metrics timestamp
	for _, metric := range metrics {
		if len(metric.Fields()) == 0 {
			s.log.Debugf("Skipping the serialization of metric %s without fields ", metric.Name())
			continue
		}

		serializedTags := serializeTags(metric.Tags())
		serializedMetrics := serializeMetrics(metric)
		serializedTimestamp := strconv.FormatInt(metric.Time().UnixNano(), 10)

		if serializedMetrics == "" {
			continue
		}

		output.WriteString(nameEscape(metric.Name()))
		if serializedTags != "" {
			output.WriteString(",")
			output.WriteString(serializedTags)
		}
		output.WriteString(" ")
		output.WriteString(serializedMetrics)
		output.WriteString(" ")
		output.WriteString(serializedTimestamp)
		// has to end with a newline
		output.WriteString("\n")
	}

	return output.Bytes()
}

func serializeTags(tags map[string]string) string {
	var serializedTags strings.Builder

	// make tag order sorted
	sortedTagKeys := make([]string, 0, len(tags))
	for t := range tags {
		sortedTagKeys = append(sortedTagKeys, t)
	}
	sort.Strings(sortedTagKeys)

	for _, tagKey := range sortedTagKeys {
		tagValue := tags[tagKey]
		if serializedTags.Len() > 0 {
			serializedTags.WriteString(",")
		}

		serializedTags.WriteString(escape(tagKey))
		serializedTags.WriteString("=")
		serializedTags.WriteString(escape(tagValue))
	}
	return serializedTags.String()
}

func serializeMetrics(metric telegraf.Metric) string {
	var serializedMetrics strings.Builder

	// make the field order sorted
	sort.Slice(metric.FieldList(), func(i, j int) bool {
		return metric.FieldList()[i].Key < metric.FieldList()[j].Key
	})

	var countAdded = 0
	for _, field := range metric.FieldList() {
		var serializedMetric = serializeMetric(field.Key, field.Value)

		if serializedMetric == "" {
			continue
		}

		if countAdded > 0 {
			serializedMetrics.WriteString(",")
		}
		serializedMetrics.WriteString(serializedMetric)
		countAdded++
	}

	return serializedMetrics.String()
}

func serializeMetric(key string, value interface{}) string {
	var metricValue string
	switch v := value.(type) {
	case string:
		metricValue = fmt.Sprintf("\"%s\"", stringFieldEscape(v))
	case bool:
		metricValue = strconv.FormatBool(v)
	case float64:
		if !math.IsNaN(v) && !math.IsInf(v, 0) {
			metricValue = strconv.FormatFloat(v, 'f', -1, 64)
		} else {
			return ""
		}
	case uint64:
		metricValue = strconv.FormatUint(v, 10) + "i"
	case int64:
		metricValue = strconv.FormatInt(v, 10) + "i"
	default:
		return ""
	}

	return fmt.Sprint(key, "=", metricValue)
}

// CompactMetricSerializer can be used to output metrics in compact format. Compact format squeezes as many metrics
// as possible in a single output line, based on tags and timestamp of each metric. When multiple metrics share the
// same tags and the same timestamp, we can write them in a single line to reduce the total bulk size by not
// repeating the same tags+timestamp multiple times.
// TODO to be implemented in the future to make Telegraf requests to Sematext backend smaller
type CompactMetricSerializer struct {
	tagsIDToMetrics map[string]telegraf.Metric
}
