package main

import (
	"github.com/mikanmekan/koalemos/cmd/doorman/server"
)

func main() {
	s := server.New(8080)
	s.HandleRequests()
}
