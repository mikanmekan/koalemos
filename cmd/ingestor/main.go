package main

import (
	"github.com/mikanmekan/koalemos/cmd/ingestor/ingestion"
	"github.com/mikanmekan/koalemos/cmd/ingestor/server"
	"github.com/mikanmekan/koalemos/internal/metrics"
)

func main() {
	reader := metrics.NewReader()
	ingestion := ingestion.New(reader, nil)
	s := server.New(8080, reader)
	s.HandleRequests()
}
