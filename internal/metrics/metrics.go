package metrics

// MetricPoint represents a single Koalemos data model metric datum.
type MetricPoint struct {
	Name      string
	Value     float64
	LabelSet  map[string]string
	Timestamp int64
}

// MetricFamily represents a group of metrics.
type MetricFamily struct {
	Name    string
	Metrics []MetricPoint
	Type    string
	Help    string
}
