package metrics

import (
	"errors"
)

var (
	ErrMetricFamilyNotFound    = errors.New("metric family not found")
	ErrTimeSeriesNotFound      = errors.New("timeseries not found")
	ErrDuplicateMetricLabelSet = errors.New("label set should not be repeated within a MetricFamiliesTimeGroup")
)
