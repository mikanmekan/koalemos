package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mikanmekan/koalemos/internal/log"
	"go.uber.org/zap"
)

// Server listens for metrics being sent by clients and ingests them.
type Server struct {
	logger log.Logger
	router *mux.Router
	port   int
}

// New initializes a Server which will listen on the given port.
func New(port int) *Server {
	s := &Server{
		logger: log.NewLogger(),
		router: mux.NewRouter(),
		port:   port,
	}

	return s
}

// HandleRequests starts the server and listens for incoming requests.
func (s *Server) HandleRequests() {
	s.router.HandleFunc("/metrics", s.handleMetrics).Methods("POST")

	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router)
	if err != nil {
		s.logger.Fatal("server failed to start serving requests", zap.Error(err))
	}
}

// handleMetrics expects a POST request with a JSON body containing metrics in
// OpenMetrics format.
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Warn("failed to read metrics", zap.Error(err))
		return
	}

	s.logger.Debug("received metrics from client", zap.String("metrics", string(metrics)))

	w.WriteHeader(http.StatusOK)
}
