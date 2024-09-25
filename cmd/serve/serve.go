package main

import (
	"log/slog"
	"os"

	a "github.com/jesses-code-adventures/slop/api"
	_ "github.com/jesses-code-adventures/slop/env"
	"github.com/jesses-code-adventures/slop/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server := server.NewServer(logger)
	api := a.NewTestApi(&server, logger)
	a.BindRoutes(api)
	server.Serve()
}
