package metrics

// Metric is a single metric in OpenMetrics format.
//
// Open Metrics data model:
// https://github.com/OpenObservability/OpenMetrics/blob/main/specification/OpenMetrics.md#data-model
type Metric struct {
	Name     string
	Value    float64
	LabelSet map[string]string
}

// Reader reads incoming byte streams for metrics.
type Reader interface {
	Read([]byte) ([]Metric, error)
}

// OMReader reads incoming byte streams for metrics in the OpenMetrics
// format.
type OMReader struct {
}
