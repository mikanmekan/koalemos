package metrics

import (
	"fmt"
	"strconv"
	"strings"
)

// MetricPoint represents a single Koalemos data model metric datum.
type MetricPoint struct {
	Name     string
	Value    float64
	LabelSet map[string]string
	Time     int64
}

func (m *MetricPoint) String() string {
	sb := strings.Builder{}

	metricPoint := "Name: " + m.Name + ", " +
		"Value: " + "%f" + ", " +
		"Name: " + m.Name + ", " +
		"Timestamp: " + "%d" + ", " +
		"Labels:"

	floatStr := strconv.FormatFloat(m.Value, 'f', 4, 64)
	metricPoint = fmt.Sprintf(metricPoint, floatStr, m.Time)

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

	metricFamily := "Name: " + m.Name + ", " +
		"Type: " + m.Type + ", " +
		"Help: " + m.Help + ", " +
		"Metrics:"

	sb.WriteString(metricFamily)
	for _, v := range m.Metrics {
		sb.WriteString(v.String())
	}

	return sb.String()
}

type MetricFamiliesTimeGroup struct {
	Time     int64
	Families map[string]*MetricFamily
}

func NewMetricFamiliesTimeGroup() *MetricFamiliesTimeGroup {
	return &MetricFamiliesTimeGroup{
		Time:     0,
		Families: map[string]*MetricFamily{},
	}
}

// Apply adds the information within mf to m.
func (m *MetricFamiliesTimeGroup) Apply(mf *MetricFamily) error {
	if mf == nil {
		return fmt.Errorf("partial metric family is nil")
	}

	if v, ok := m.Families[mf.Name]; ok {
		// apply non-zero values
		if len(mf.Metrics) == 0 {
			v.Metrics = mf.Metrics
		}
		if mf.Help != "" {
			v.Help = mf.Help
		}
		if mf.Type != "" {
			v.Type = mf.Type
		}
	} else {
		m.Families[mf.Name] = mf
	}
	return nil
}
