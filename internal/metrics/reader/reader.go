package reader

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unsafe"

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
	Read(io.Reader) (map[string]*metrics.MetricFamily, error)
}

// MetricsReader reads incoming byte streams for metrics in the Koalemos
// format.
type MetricsReader struct{}

var _ Reader = (*MetricsReader)(nil)

func NewReader() *MetricsReader {
	return &MetricsReader{}
}

// Read incoming byte streams for metrics in the Koalemos format.
func (r *MetricsReader) Read(requestReader io.Reader) (map[string]*metrics.MetricFamily, error) {
	metricFamilies := make(map[string]*metrics.MetricFamily)

	scanner := bufio.NewScanner(requestReader)
	for scanner.Scan() {
		line := BytesToString(scanner.Bytes())
		err := processLine(line, metricFamilies)
		if err != nil {
			return metricFamilies, err
		}
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return metricFamilies, nil
}

func BytesToString(b []byte) string {
	p := unsafe.SliceData(b)
	return unsafe.String(p, len(b))
}

// processLine takes a byte slice representing a line in the metrics payload,
// and applies the information to the metricFamilies.
func processLine(line string, metricFamilies map[string]*metrics.MetricFamily) error {
	var err error
	if line[0] == '#' {
		err = stripMetricFamilyMetadata(line, metricFamilies)
	} else {
		err = processMetric(line)
	}
	return err
}

func processMetric(line string) error {
	return nil
}

func stripMetricFamilyMetadata(line string, metricFamilies map[string]*metrics.MetricFamily) error {
	// metadataString is the full line `# HELP metric desc....` ->
	// `HELP metric desc....`
	metadataString := line[2:]
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
		fmt.Println(metadataPieces, metadataString)
		fmt.Println("----------")
		return ErrUnexpectedMetadata
	}

	return nil
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
