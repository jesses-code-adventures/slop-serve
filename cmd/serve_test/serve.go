package main

import (
	"context"
	"log/slog"
	"os"

	a "github.com/jesses-code-adventures/slop/api"
	_ "github.com/jesses-code-adventures/slop/env"
	"github.com/jesses-code-adventures/slop/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx := context.WithValue(context.Background(), "logger", logger)
	server := server.NewServer(ctx)
	api := a.NewTestApi(&server, logger)
	a.BindRoutes(api)
	logger.Info("running serve...")
	server.Serve()
}
