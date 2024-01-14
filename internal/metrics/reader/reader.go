package reader

import (
	"fmt"
	"strings"

	"github.com/mikanmekan/koalemos/internal/metrics"
)

type metricComponent int

const (
	TYPE metricComponent = iota
	NAME
	TEXT
)

// Reader reads incoming byte streams for metrics.
type Reader interface {
	Read([]byte) (map[string]*metrics.MetricFamily, error)
}

// MetricsReader reads incoming byte streams for metrics in the Koalemos
// format.
type MetricsReader struct{}

var _ Reader = (*MetricsReader)(nil)

func NewReader() *MetricsReader {
	return &MetricsReader{}
}

// Read incoming byte streams for metrics in the Koalemos format.
func (r *MetricsReader) Read(bytes []byte) (map[string]*metrics.MetricFamily, error) {
	metricFamilies := make(map[string]*metrics.MetricFamily)

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

func processLine(bytes []byte, metricFamilies map[string]*metrics.MetricFamily, i int) (int, error) {
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

func stripMetricFamilyMetadata(bytes []byte, metricFamilies map[string]*metrics.MetricFamily, i int) (int, error) {
	end := i
	for end = i; end < len(bytes); end++ {
		if bytes[end] == '\n' {
			break
		}
	}
	// metadataString is the full line `# HELP metric desc....` ->
	// `HELP metric desc....`
	metadataString := string(bytes[i+1 : end])
	// Split such that the first element is the metric family name, and the
	// second is the relevant metadata.
	metadataPieces := strings.SplitN(metadataString, " ", 3)

	switch metadataPieces[TYPE] {
	case "TYPE":
		enrichMetricFamilies(metricFamilies, &metrics.MetricFamily{
			Name: metadataPieces[NAME],
			Type: "gauge", // TO-DO: Support other metric types.
		})
	case "HELP":
		// Assuming we encounter HELP before any other metadata or metrics.
		enrichMetricFamilies(metricFamilies, &metrics.MetricFamily{
			Name: metadataPieces[NAME],
			Help: metadataPieces[TEXT],
		})
	default:
		fmt.Println("encountered unexpected metadata:")
		fmt.Println(i, metadataPieces, metadataString)
		fmt.Println("----------")
		return i + 1, ErrUnexpectedMetadata
	}

	return end, nil
}

// enrichMetricFamilies takes a metric family and adds the input metric family's
// information to the metricFamilies.
func enrichMetricFamilies(metricFamilies map[string]*metrics.MetricFamily, partialMetricFamily *metrics.MetricFamily) error {
	var (
		metricFamilyPtr *metrics.MetricFamily
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
