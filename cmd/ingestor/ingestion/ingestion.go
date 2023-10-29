package ingestion

import (
	"github.com/mikanmekan/koalemos/internal/metrics"
	"github.com/mikanmekan/koalemos/internal/metrics/store"
)

func New(mr metrics.Reader, ims store.IMS) *Ingestor {
	return &Ingestor{
		metricsReader: mr,
		metricsIMS:    ims,
	}
}

type Ingestor struct {
	metricsReader metrics.Reader
	metricsIMS    store.IMS
}
