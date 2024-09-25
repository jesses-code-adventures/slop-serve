package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	_ "github.com/jesses-code-adventures/slop/env"
)

type Server struct {
	Port string
	Mux  *http.ServeMux
}

func NewServer() Server {
	port := os.Getenv("DEV_PORT")
	if len(port) == 0 {
		panic("no port found")
	}
	if !strings.HasPrefix(port, ":") {
		port = fmt.Sprintf(":%s", port)
	}
	mux := http.NewServeMux()
	return Server{
		Port: port,
		Mux:  mux,
	}
}

func (s *Server) Serve() {
	err := http.ListenAndServe(s.Port, s.Mux)
	if err != nil {
		panic(err)
	}
}
