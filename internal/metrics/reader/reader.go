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
	hashstructure "github.com/mitchellh/hashstructure/v2"
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
			return fmt.Errorf("processing first line: %w", err)
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

// processMetric with input line in format metric_name{lbl1="val",lbl2="val"} 10
func processMetric(line string, metricFamilies *metrics.MetricFamiliesTimeGroup) error {
	const (
		NAME_PART = iota
		LABEL_PART
	)

	lineParts := strings.SplitN(line, "{", -1)

	if len(lineParts) != 2 {
		return ErrUnexpectedMetricLine
	}

	labelSetParts := labelsetRegex.FindAllStringSubmatch(lineParts[LABEL_PART], -1)

	mp := metrics.MetricPoint{
		Name:     lineParts[NAME_PART],
		LabelSet: map[string]string{},
	}

	// Take hash of metric name + label set
	hash, err := hashstructure.Hash(mp, hashstructure.FormatV2, nil)
	if err != nil {
		return fmt.Errorf("hashing metric point: %w", err)
	}
	mp.Hash = hash

	err = processLabelSets(&mp, labelSetParts)
	if err != nil {
		return fmt.Errorf("processing label sets: %w", err)
	}

	val, err := parseValue(line)
	if err != nil {
		return fmt.Errorf("parsing metric value: %w", err)
	}
	mp.Value = val

	err = metricFamilies.AddMetricPoint(&mp)
	if err != nil {
		return fmt.Errorf("adding metric point: %w", err)
	}

	return nil
}

func parseValue(line string) (float64, error) {
	valueStr := strings.Fields(line)

	if len(valueStr) != 2 {
		return 0, ErrInvalidValue
	}

	value, err := strconv.ParseFloat(valueStr[1], 64)
	if err != nil {
		return 0, fmt.Errorf("string field was not a valid float value: %w", err)
	}

	return value, nil
}

func processLabelSets(mp *metrics.MetricPoint, labelSetParts [][]string) error {
	const (
		KEY   = 1
		VALUE = 2
	)

	for i := 0; i < len(labelSetParts); i++ {
		// Return err if there's a repeated key.
		if _, found := mp.LabelSet[labelSetParts[i][KEY]]; !found {
			mp.LabelSet[labelSetParts[i][KEY]] = labelSetParts[i][VALUE]
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
		def := metrics.MetricDefinition{
			Name: metadataPieces[NAME],
			Type: "gauge", // TO-DO: Support other metric types.
		}
		m := metrics.NewMetricFamily(def)
		metricFamilies.AddMetricFamily(&m)
	case "HELP":
		def := metrics.MetricDefinition{
			Name: metadataPieces[NAME],
			Help: metadataPieces[TEXT],
		}
		m := metrics.NewMetricFamily(def)
		metricFamilies.AddMetricFamily(&m)
	default:
		return ErrUnexpectedMetadata
	}

	return nil
}
