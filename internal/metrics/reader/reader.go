package reader

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"github.com/mikanmekan/koalemos/internal/metrics"
)

type metricComponent int

const (
	TYPE metricComponent = iota
	NAME
	TEXT

	// Match valid labelsets {method="post",code="400"}
	// https://regex101.com/r/IgomLp/1
	LABELSET_REGEX = `([a-zA-Z_][a-zA-Z0-9_]*?)="([a-zA-Z0-9_]*?)"(,|})`
)

var labelsetRegex = regexp.MustCompile(LABELSET_REGEX)

// Reader reads incoming byte streams for metrics.
type Reader interface {
	Read(io.Reader) (*metrics.MetricFamiliesTimeGroup, error)
}

// MetricsReader reads incoming byte streams for metrics in the Koalemos
// format.
type MetricsReader struct{}

var _ Reader = (*MetricsReader)(nil)

func NewReader() *MetricsReader {
	return &MetricsReader{}
}

// Read incoming byte streams for metrics in the Koalemos format.
func (r *MetricsReader) Read(requestReader io.Reader) (*metrics.MetricFamiliesTimeGroup, error) {
	metricFamilies := metrics.NewMetricFamiliesTimeGroup()

	scanner := bufio.NewScanner(requestReader)
	// Grab timestamp for metrics payload.
	if scanner.Scan() {
		line := BytesToString(scanner.Bytes())
		processFirstLine(line, metricFamilies)
	}
	// Process through rest of metrics payload (metadata, metrics).
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

func processFirstLine(line string, metricFamilies *metrics.MetricFamiliesTimeGroup) error {
	time, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		err := processLine(line, metricFamilies)
		if err != nil {
			return fmt.Errorf("failed to process first line: %w", err)
		}
	}

	metricFamilies.Time = time
	return nil
}

func BytesToString(b []byte) string {
	p := unsafe.SliceData(b)
	return unsafe.String(p, len(b))
}

// processLine takes a byte slice representing a line in the metrics payload,
// and applies the information to the metricFamilies.
func processLine(line string, metricFamilies *metrics.MetricFamiliesTimeGroup) error {
	var err error
	if line[0] == '#' {
		err = stripMetricFamilyMetadata(line, metricFamilies)
	} else {
		err = processMetric(line, metricFamilies)
	}
	return err
}

func processMetric(line string, metricFamilies *metrics.MetricFamiliesTimeGroup) error {
	lineParts := strings.SplitN(line, "{", 2)

	if len(lineParts) != 2 {
		return ErrUnexpectedMetricLine
	}

	labelSetParts := labelsetRegex.FindAllStringSubmatch(lineParts[1], -1)

	mp := metrics.MetricPoint{
		Name:     lineParts[0],
		LabelSet: map[string]string{},
	}

	err := processLabelSets(&mp, labelSetParts)
	if err != nil {
		return err
	}

	metricFamilies.AddMetricPoint(&mp)

	return nil
}

func processLabelSets(mp *metrics.MetricPoint, labelSetParts [][]string) error {
	if len(labelSetParts)%2 == 1 {
		return ErrOddLabelSetParts
	}

	for i := 0; i < len(labelSetParts); i++ {
		// Return err if there's a repeated key.
		if _, found := mp.LabelSet[labelSetParts[i][1]]; !found {
			mp.LabelSet[labelSetParts[i][1]] = labelSetParts[i][2]
		} else {
			return ErrDuplicateLabelKey
		}
	}

	return nil
}

func stripMetricFamilyMetadata(line string, metricFamilies *metrics.MetricFamiliesTimeGroup) error {
	// metadataString is the full line `# HELP metric desc....` ->
	// `HELP metric desc....`
	metadataString := line[2:]
	// Split such that the first element is the metric family name, and the
	// second is the relevant metadata.
	metadataPieces := strings.SplitN(metadataString, " ", 3)

	switch metadataPieces[TYPE] {
	case "TYPE":
		metricFamilies.AddMetricFamily(&metrics.MetricFamily{
			Name: metadataPieces[NAME],
			Type: "gauge", // TO-DO: Support other metric types.
		})
	case "HELP":
		// Assuming we encounter HELP before any other metadata or metrics.
		metricFamilies.AddMetricFamily(&metrics.MetricFamily{
			Name: metadataPieces[NAME],
			Help: metadataPieces[TEXT],
		})
	default:
		return ErrUnexpectedMetadata
	}

	return nil
}
