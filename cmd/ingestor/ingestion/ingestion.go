package ingestion

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mikanmekan/koalemos/internal/log"
	reader "github.com/mikanmekan/koalemos/internal/metrics/reader"
	"github.com/mikanmekan/koalemos/internal/metrics/store"
	"go.uber.org/zap"
)

func New(l log.Logger, mr reader.Reader, ims store.IMS) *Ingestor {
	return &Ingestor{
		logger:     l,
		metricsIMS: ims,
	}
}

type Ingestor struct {
	logger     log.Logger
	metricsIMS store.IMS
}

// HandleMetrics expects a POST request with a JSON body containing metrics in
// Koalemos format.
func (i *Ingestor) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	metricsReader := reader.NewReader()

	mfs, err := metricsReader.Read(r.Body)
	if err != nil {
		i.logger.Warn("failed to read metrics", zap.Error(err))
		return
	}

	i.logger.Info(fmt.Sprintf("%+v", mfs))

	err = i.metricsIMS.AddMetricFamiliesTimeGroup(mfs)
	if err != nil {
		i.logger.Error("failed to write metrics to in memory store", zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (i *Ingestor) Register(r *mux.Router) {
	r.HandleFunc("/metrics", i.HandleMetrics).Methods("POST")
}
