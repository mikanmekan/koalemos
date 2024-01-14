package ingestion

import (
	"fmt"
	"io"
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

	metrics, err := io.ReadAll(r.Body)
	if err != nil {
		i.logger.Warn("failed to read metrics", zap.Error(err))
		return
	}

	mfs, err := metricsReader.Read(metrics)
	if err != nil {
		i.logger.Warn("failed to read metrics", zap.Error(err))
		return
	}

	i.logger.Info(fmt.Sprintf("%v", mfs))

	w.WriteHeader(http.StatusOK)
}

func (i *Ingestor) Register(r *mux.Router) {
	r.HandleFunc("/metrics", i.HandleMetrics).Methods("POST")
}
