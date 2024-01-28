package metrics

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// MetricPoint represents a single Koalemos data model metric datum.
type MetricPoint struct {
	Name     string
	Value    float64
	LabelSet map[string]string
	Time     int64
	Hash     uint64
}

// MetadataEquals will return true between m1 & m2 if the two metric points'
// label sets are equal. This is used to check whether two metric points can
// both reside within a MetricFamiliesTimeGroup.
func (m1 *MetricPoint) MetadataEquals(m2 *MetricPoint) bool {
	if m1.Hash != m2.Hash {
		return false
	}

	if !reflect.DeepEqual(m1.LabelSet, m2.LabelSet) {
		return false
	}

	return true
}

func (m *MetricPoint) String() string {
	sb := strings.Builder{}

	metricPoint := "Name: " + m.Name + ", " +
		"Value: " + "%s" + ", " +
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

// MetricDefinition uniquely defines a metric.
type MetricDefinition struct {
	Name string
	Type string
	Help string
}

// MetricFamily represents a group of metrics.
type MetricFamily struct {
	Definition    MetricDefinition
	HashedMetrics map[uint64][]*MetricPoint
}

func NewMetricFamily(def MetricDefinition) MetricFamily {
	m := MetricFamily{HashedMetrics: map[uint64][]*MetricPoint{}}
	m.Definition = def
	return m
}

func (m *MetricFamily) String() string {
	sb := strings.Builder{}

	metricFamily := "Name: " + m.Definition.Name + ", " +
		"Type: " + m.Definition.Type + ", " +
		"Help: " + m.Definition.Help + ", " +
		"Metrics:"

	sb.WriteString(metricFamily)
	for _, hashSlice := range m.HashedMetrics {
		for _, mp := range hashSlice {
			sb.WriteString(mp.String())
		}
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

// AddMetricFamily adds the information within mf to m. Adds will blindly overwrite
// pre-existing information in m if present.
func (m *MetricFamiliesTimeGroup) AddMetricFamily(mf *MetricFamily) error {
	if mf == nil {
		return fmt.Errorf("partial metric family is nil")
	}

	if v, found := m.Families[mf.Definition.Name]; found {
		// apply non-zero values
		if len(mf.HashedMetrics) == 0 {
			v.HashedMetrics = mf.HashedMetrics
		}
		if mf.Definition.Help != "" {
			v.Definition.Help = mf.Definition.Help
		}
		if mf.Definition.Type != "" {
			v.Definition.Type = mf.Definition.Type
		}
	} else {
		m.Families[mf.Definition.Name] = mf
	}
	return nil
}

func (m *MetricFamiliesTimeGroup) AddMetricPoint(mp *MetricPoint) error {
	mf, err := m.GetMetricFamily(mp.Name)
	if err != nil {
		return fmt.Errorf("adding metric point: %w", err)
	}

	// Only add metric if there is no pre-existing labelset which would be
	// colliding with this metric point.
	if err := checkCollision(mf, mp); err != nil {
		return fmt.Errorf("adding metric point: %w", err)
	}

	// Add metric, indexed by hash
	m.Families[mp.Name].HashedMetrics[mp.Hash] = append(m.Families[mp.Name].HashedMetrics[mp.Hash], mp)

	return nil
}

func (m *MetricFamiliesTimeGroup) GetMetricFamily(metricName string) (*MetricFamily, error) {
	if v, ok := m.Families[metricName]; ok {
		return v, nil
	}
	return nil, ErrMetricFamilyNotFound
}

// checkCollision returns ErrDuplicateMetricLabelSet if mfs contains
// mp. Also update metric family inside mfs to contain the new hash.
func checkCollision(mf *MetricFamily, mp *MetricPoint) error {
	hashedMps, found := mf.HashedMetrics[mp.Hash]
	if !found {
		return nil
	}

	for _, hashedMp := range hashedMps {
		if hashedMp.MetadataEquals(mp) {
			return ErrDuplicateMetricLabelSet
		}
	}

	// (to-do: return true if mp matches value found above)
	return nil
}

// MetricFamilyTimeSeries
type MetricFamilyTimeSeries struct {
	Def      MetricDefinition
	LabelSet map[string]string
	metrics  []MetricPoint
}
