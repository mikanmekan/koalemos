package store

import "github.com/mikanmekan/koalemos/internal/metrics"

// MetricsIMS is the interface for an in memory store for metrics.
type IMS interface {
	AddMetricFamiliesTimeGroup(metricFamiliesTimeGroup *metrics.MetricFamiliesTimeGroup) error
	GetMetricFamily(metricFamily metrics.MetricFamily) error
}

// MetricsIMSImpl is the in memory store for metrics.
type IMSImpl struct{}

var _ IMS = (*IMSImpl)(nil)

func New() *IMSImpl {
	return &IMSImpl{}
}

func (ims *IMSImpl) AddMetricFamiliesTimeGroup(metricFamiliesTimeGroup *metrics.MetricFamiliesTimeGroup) error {
	// take each metric family
	//

	// for i, metricFamily := range metricFamiliesTimeGroup.Families {
	// 	metricFamily.Metrics
	// }
}

func (ims *IMSImpl) GetMetricFamily(metricFamily metrics.MetricFamily) error {
	panic("not implemented")
}
