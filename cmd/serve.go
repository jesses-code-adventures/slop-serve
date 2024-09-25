package main

import (
	_ "github.com/jesses-code-adventures/slop/env"
	"github.com/jesses-code-adventures/slop/server"
)

func main() {
	server := server.NewServer()
	server.Serve()
}
