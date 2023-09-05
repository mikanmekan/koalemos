package main

import (
	"github.com/mikanmekan/koalemos/server"
)

func main() {
	s := server.New(8080)
	s.HandleRequests()
}
