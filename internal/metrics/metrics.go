package metrics

import (
	"fmt"
	"strings"
)

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
type MetricsReader struct{}

var _ Reader = (*MetricsReader)(nil)

func NewReader() *MetricsReader {
	return &MetricsReader{}
}

// Read incoming byte streams for metrics in the Koalemos format.
func (r *MetricsReader) Read(bytes []byte) (map[string]*MetricFamily, error) {
	metricFamilies := make(map[string]*MetricFamily)

	// TO-DO: Benchmark this vs alloc heavy strings.Split()s.
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
	metadataString := string(bytes[i:end])
	// Split such that the first element is the metric family name, and the
	// second is the relevant metadata.
	metadataPieces := strings.SplitN(metadataString, " ", 1)

	switch metadataPieces[0] {
	case "TYPE":
		enrichMetricFamilies(metricFamilies, &MetricFamily{
			Name: metadataPieces[0],
			Type: "gauge", // TO-DO: Support other metric types.
		})
	case "HELP":
		// Assuming we encounter HELP before any other metadata or metrics.
		enrichMetricFamilies(metricFamilies, &MetricFamily{
			Name: metadataPieces[0],
			Help: metadataPieces[1],
		})
	default:
		fmt.Println(i, metadataPieces, metadataString)
		return i + 1, ErrUnexpectedMetadata
	}

	return i, nil
}

// enrichMetricFamilies takes a metric family and adds the input metric family's
// information to the metricFamilies.
func enrichMetricFamilies(metricFamilies map[string]*MetricFamily, partialMetricFamily *MetricFamily) error {
	var (
		metricFamilyPtr *MetricFamily
		ok              bool
	)

	if partialMetricFamily == nil {
		return fmt.Errorf("partial metric family is nil")
	}

	// If this is the first insertion of data, just add the pointer to the map.
	if metricFamilyPtr, ok = metricFamilies[partialMetricFamily.Name]; !ok {
		metricFamilies[partialMetricFamily.Name] = partialMetricFamily
		return nil
	}

	// If the metric family already exists, enrich existing data.
	switch {
	case partialMetricFamily.Help != "":
		metricFamilyPtr.Help = partialMetricFamily.Help
	case partialMetricFamily.Type != "":
		metricFamilyPtr.Type = partialMetricFamily.Type
	}

	return nil
}
