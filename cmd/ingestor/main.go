package main

import (
	"github.com/mikanmekan/koalemos/cmd/ingestor/ingestion"
	"github.com/mikanmekan/koalemos/cmd/ingestor/server"
	"github.com/mikanmekan/koalemos/internal/log"
	"github.com/mikanmekan/koalemos/internal/metrics"
)

func main() {
	reader := metrics.NewReader()
	logger := log.NewLogger()
	ingestion := ingestion.New(logger, reader, nil)
	s := server.New(8080, *ingestion)
	s.HandleRequests()
}
