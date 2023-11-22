package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mikanmekan/koalemos/cmd/ingestor/ingestion"
	"github.com/mikanmekan/koalemos/internal/log"
	"go.uber.org/zap"
)

// Server listens for metrics being sent by clients and ingests them.
type Server struct {
	logger   log.Logger
	router   *mux.Router
	port     int
	ingestor ingestion.Ingestor
}

// New initializes a Server which will listen on the given port.
func New(port int, ingestor ingestion.Ingestor) *Server {
	s := &Server{
		logger:   log.NewLogger(),
		router:   mux.NewRouter(),
		port:     port,
		ingestor: ingestor,
	}

	return s
}

// HandleRequests starts the server and listens for incoming requests.
func (s *Server) HandleRequests() {
	s.ingestor.Register(s.router)

	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router)
	if err != nil {
		s.logger.Fatal("server failed to start serving requests", zap.Error(err))
	}
}
