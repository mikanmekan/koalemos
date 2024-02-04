package store

import (
	"github.com/mikanmekan/koalemos/internal/hash"
	"github.com/mikanmekan/koalemos/internal/metrics"
)

// MetricsIMS is the interface for an in memory store for metrics.
type IMS interface {
	AddMetricFamiliesTimeGroup(metricFamiliesTimeGroup *metrics.MetricFamiliesTimeGroup) error
	GetTimeSeries(timeSeries *metrics.MetricFamilyTimeSeries) (*metrics.MetricFamilyTimeSeries, error)
}

// MetricsIMSImpl is the in memory store for metrics.
type IMSImpl struct {
	// activeBlock is the metrics block that all incoming metrics will be written to.
	activeBlock metrics.Block
}

var _ IMS = (*IMSImpl)(nil)

func New() *IMSImpl {
	return &IMSImpl{}
}

// AddMetricFamiliesTimeGroup adds all metrics read in from a metrics payload.
func (ims *IMSImpl) AddMetricFamiliesTimeGroup(metricFamiliesTimeGroup *metrics.MetricFamiliesTimeGroup) error {
	for metricName, metricFamily := range metricFamiliesTimeGroup.Families {
		hashed := hash.HashString(metricName)
		ims.activeBlock.AddMetricFamily(metricFamily, hashed)
	}
	return nil
}

func (ims *IMSImpl) GetTimeSeries(timeSeries *metrics.MetricFamilyTimeSeries) (*metrics.MetricFamilyTimeSeries, error) {
	return nil, nil
}
