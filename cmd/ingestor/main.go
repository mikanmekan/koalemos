package main

import (
	"github.com/mikanmekan/koalemos/cmd/ingestor/server"
)

func main() {
	s := server.New(8080)
	s.HandleRequests()
}
