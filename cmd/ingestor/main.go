package main

import (
	"github.com/mikanmekan/koalemos/cmd/ingestor/ingestion"
	"github.com/mikanmekan/koalemos/cmd/ingestor/server"
)

func main() {
	ingestion := ingestion.New()
	s := server.New(8080)
	s.HandleRequests()
}
