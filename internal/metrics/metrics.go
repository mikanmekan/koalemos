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
		"Value: " + "%f" + ", " +
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

// MetricFamily represents a group of metrics.
type MetricFamily struct {
	Name    string
	Metrics []MetricPoint
	Type    string
	Help    string
	Hashes  map[uint64][]*MetricPoint
}

func NewMetricFamily() MetricFamily {
	return MetricFamily{Hashes: map[uint64][]*MetricPoint{}}
}

func (m *MetricFamily) String() string {
	sb := strings.Builder{}

	metricFamily := "Name: " + m.Name + ", " +
		"Type: " + m.Type + ", " +
		"Help: " + m.Help + ", " +
		"Metrics:"

	sb.WriteString(metricFamily)
	for _, v := range m.Metrics {
		sb.WriteString(v.String())
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

	if v, found := m.Families[mf.Name]; found {
		// apply non-zero values
		if len(mf.Metrics) == 0 {
			v.Metrics = mf.Metrics
		}
		if mf.Help != "" {
			v.Help = mf.Help
		}
		if mf.Type != "" {
			v.Type = mf.Type
		}
	} else {
		m.Families[mf.Name] = mf
	}
	return nil
}

func (m *MetricFamiliesTimeGroup) AddMetricPoint(mp *MetricPoint) error {
	// Only add metric if there is no pre-existing labelset which would be
	// colliding with this metric point.
	if err := checkCollision(m, mp); err != nil {
		return fmt.Errorf("adding metric point: %w", err)
	}

	// To-do - I wrote metric points as structs not pointer to structs for
	// locality. Do we actually need to copy these structs?
	m.Families[mp.Name].Metrics = append(m.Families[mp.Name].Metrics, *mp)

	// Update hash - (to-do: Add method for this, this is an eyesore!)
	m.Families[mp.Name].Hashes[mp.Hash] = append(m.Families[mp.Name].Hashes[mp.Hash], mp)

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
func checkCollision(mfs *MetricFamiliesTimeGroup, mp *MetricPoint) error {
	mf, err := mfs.GetMetricFamily(mp.Name)
	if err != nil {
		return nil
	}

	hashedMps, found := mf.Hashes[mp.Hash]
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
