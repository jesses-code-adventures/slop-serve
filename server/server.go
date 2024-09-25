package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	_ "github.com/jesses-code-adventures/slop/env"
)

type Server struct {
	Port   string
	Mux    *http.ServeMux
	Logger *slog.Logger
}

func NewServer(logger *slog.Logger) Server {
	port := os.Getenv("SERVER_PORT")
	if len(port) == 0 {
		panic("no port found")
	}
	if !strings.HasPrefix(port, ":") {
		port = fmt.Sprintf(":%s", port)
	}
	mux := http.NewServeMux()
	return Server{
		Port:   port,
		Mux:    mux,
		Logger: logger,
	}
}

func (s *Server) Serve() {
	err := http.ListenAndServe(s.Port, s.Mux)
	if err != nil {
		panic(err)
	}
}

func (s *Server) RegisterHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.Mux.HandleFunc(pattern, handler)
}
