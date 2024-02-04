package metrics

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/hashstructure/v2"
)

// MetricPoint represents a single Koalemos data model metric datum.
type MetricPoint struct {
	Name     string
	Value    float64 `hash:"ignore"`
	LabelSet map[string]string
	Time     int64  `hash:"ignore"`
	Hash     uint64 `hash:"ignore"`
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

func (mf *MetricFamily) ToMetricTimeSeries() *MetricFamilyTimeSeries {
	return &MetricFamilyTimeSeries{
		Def: mf.Def,
	}
}

// HashMetric takes a hash of metric name + label set and returns the hash.
func HashMetric(mp *MetricPoint) (uint64, error) {
	hash, err := hashstructure.Hash(mp, hashstructure.FormatV2, nil)
	if err != nil {
		return 0, fmt.Errorf("hashing metric point: %w", err)
	}
	return hash, nil
}

// MetadataEquals will return true between m1 & m2 if the two metric points'
// label sets are equal. This is used to check whether two metric points can
// both reside within a MetricFamiliesTimeGroup.
func (m1 *MetricPoint) MetadataEquals(m2 *MetricPoint) bool {
	if m1.Hash != m2.Hash {
		return false
	}

	if m1.Name != m2.Name {
		return false
	}

	if !reflect.DeepEqual(m1.LabelSet, m2.LabelSet) {
		return false
	}

	return true
}

// MetricDefinition uniquely defines a metric.
type MetricDefinition struct {
	Name string
	Type string
	Help string
}

// MetricFamily represents a group of metrics.
type MetricFamily struct {
	Def           MetricDefinition
	HashedMetrics map[uint64][]*MetricPoint
}

func NewMetricFamily(def MetricDefinition) MetricFamily {
	m := MetricFamily{HashedMetrics: map[uint64][]*MetricPoint{}}
	m.Def = def
	return m
}

func (m *MetricFamily) String() string {
	sb := strings.Builder{}

	metricFamily := "Name: " + m.Def.Name + ", " +
		"Type: " + m.Def.Type + ", " +
		"Help: " + m.Def.Help + ", " +
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

	if v, found := m.Families[mf.Def.Name]; found {
		// apply non-zero values
		if len(mf.HashedMetrics) == 0 {
			v.HashedMetrics = mf.HashedMetrics
		}
		if mf.Def.Help != "" {
			v.Def.Help = mf.Def.Help
		}
		if mf.Def.Type != "" {
			v.Def.Type = mf.Def.Type
		}
	} else {
		m.Families[mf.Def.Name] = mf
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

// checkCollision returns ErrDuplicateMetricLabelSet if mf contains
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

	return nil
}

// MetricFamilyTimeSeries
type MetricFamilyTimeSeries struct {
	Def      MetricDefinition
	LabelSet map[string]string
	metrics  []MetricPoint
}

func (ts *MetricFamilyTimeSeries) Hash() (uint64, error) {
	mp := MetricPoint{
		Name:     ts.Def.Name,
		LabelSet: ts.LabelSet,
	}
	return HashMetric(&mp)
}

func (ts *MetricFamilyTimeSeries) ToMetricPoint() *MetricPoint {
	return &MetricPoint{
		Name:     ts.Def.Name,
		LabelSet: ts.LabelSet,
	}
}
