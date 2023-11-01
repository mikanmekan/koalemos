package metrics

// MetricPoint represents a single Koalemos data model metric datum.
type MetricPoint struct {
	Name      string
	Value     float64
	LabelSet  map[string]string
	Timestamp int64
}

// Reader reads incoming byte streams for metrics.
type Reader interface {
	Read([]byte) ([]MetricPoint, error)
}

// MetricsReader reads incoming byte streams for metrics in the Koalemos
// format.
type MetricsReader struct {
}

func New() *MetricsReader {
	return &MetricsReader{}
}

func (r *MetricsReader) Read(bytes []byte) ([]MetricPoint, error) {
	metrics := make([]MetricPoint, 0)
	return metrics, nil
}
