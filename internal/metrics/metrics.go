package metrics

import "strings"

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

// Reader reads incoming byte streams for metrics.
type Reader interface {
	Read([]byte) (map[string]*MetricFamily, error)
}

// MetricsReader reads incoming byte streams for metrics in the Koalemos
// format.
type MetricsReader struct {
}

var _ Reader = (*MetricsReader)(nil)

func New() *MetricsReader {
	return &MetricsReader{}
}

// Read incoming byte streams for metrics in the Koalemos format.
func (r *MetricsReader) Read(bytes []byte) (map[string]*MetricFamily, error) {
	metricFamilies := make(map[string]*MetricFamily)

	// TO-DO (p2): Benchmark this vs alloc heavy strings.Split()s.
	// Split the incoming byte stream into individual metrics.
	var err error
	for i := 0; i < len(bytes); i++ {
		i, err = processLine(bytes, metricFamilies, i)
		if err != nil {
			return metricFamilies, err
		}
	}

	return metricFamilies, nil
}

func processLine(bytes []byte, metricFamilies map[string]*MetricFamily, i int) (int, error) {
	var err error
	if bytes[i] == '#' {
		i, err = stripMetricFamilyMetadata(bytes, metricFamilies, i+1)
	} else {
		i, err = processMetric(bytes, i)
	}
	return i, err
}

func processMetric(bytes []byte, i int) (int, error) {
	return i, nil
}

func stripMetricFamilyMetadata(bytes []byte, metricFamilies map[string]*MetricFamily, i int) (int, error) {
	end := i
	for end = i; end < len(bytes); end++ {
		if bytes[i] == '\n' {
			break
		}
	}
	mdString := string(bytes[i:end])
	metadataPieces := strings.SplitN(mdString, " ", 1)

	switch metadataPieces[0] {
	case "TYPE":
		// Discard for now, only supporting ambiguous float type.
	case "HELP":
		// Assuming we encounter HELP before any other metadata or metrics.
		metricFamilies[metadataPieces[0]] = &MetricFamily{
			Name: metadataPieces[0],
			Type: "gauge",
			Help: metadataPieces[1],
		}
	default:
		return i + 1, ErrUnexpectedMetadata
	}

	return i, nil
}