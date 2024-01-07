package metrics

import (
	"fmt"
	"strconv"
	"strings"
)

// MetricPoint represents a single Koalemos data model metric datum.
type MetricPoint struct {
	Name      string
	Value     float64
	LabelSet  map[string]string
	Timestamp int64
}

func (m *MetricPoint) String() string {
	sb := strings.Builder{}

	metricPoint := "Name: " + m.Name + "\n" +
		"Value: " + "%f" + "\n" +
		"Name: " + m.Name + "\n" +
		"Timestamp: " + "%d" + "\n" +
		"Labels:"

	floatStr := strconv.FormatFloat(m.Value, 'f', 4, 64)
	metricPoint = fmt.Sprintf(metricPoint, floatStr, m.Timestamp)

	sb.WriteString(metricPoint)
	for k, v := range m.LabelSet {
		sb.WriteString(fmt.Sprintf(" {%s: %s}", k, v))
	}

	return sb.String()
}

// MetricFamily represents a group of metrics.
type MetricFamily struct {
	Name    string
	Metrics []MetricPoint
	Type    string
	Help    string
}

func (m *MetricFamily) String() string {
	sb := strings.Builder{}

	metricFamily := "Name: " + m.Name + "\n" +
		"Type: " + m.Type + "\n" +
		"Help: " + m.Help + "\n" +
		"Metrics:"

	sb.WriteString(metricFamily)
	for _, v := range m.Metrics {
		sb.WriteString(v.String())
	}

	return sb.String()
}
