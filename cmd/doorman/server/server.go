package server

import (
	"bytes"
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
		s.logger.Warn("server failed to start serving requests", zap.Error(err))
		return
	}

	fmt.Printf("Received metrics from client, request body: \n%s \n", metrics)

	// Send metrics payload onto a valid ingestor.
	s.forwardMetricsPayload(metrics)

	w.WriteHeader(http.StatusOK)
}

type Endpoint struct {
	Address string
}

// forwardMetricsPayload forwards the metrics payload to an ingestor.
// The ingestor is picked from a pool of available ingestors, where if
// the client has already been assigned a target ingestor, it will be used.
// If not, a new ingestor will be assigned to the client.
func (s *Server) forwardMetricsPayload(metrics []byte) {
	ingestors := []Endpoint{{Address: "localhost:8080"}}
	s.sendMetricsToIngestor(metrics, ingestors[0])
}

func (s *Server) sendMetricsToIngestor(metrics []byte, ingestor Endpoint) {
	// Send metrics payload to ingestor.
	_, err := http.Post(ingestor.Address, "application/json", bytes.NewBuffer(metrics))
	if err != nil {
		s.logger.Warn("failed to send metrics to ingestor", zap.Error(err))
		return
	}
}
