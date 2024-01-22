package metrics

import (
	"errors"
)

var (
	ErrMetricFamilyNotFound    = errors.New("metric family not found")
	ErrDuplicateMetricLabelSet = errors.New("label set should not be repeated within a MetricFamiliesTimeGroup")
)
