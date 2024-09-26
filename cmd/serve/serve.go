package main

import (
	"context"
	"log/slog"
	"os"

	a "github.com/jesses-code-adventures/slop/api"
	"github.com/jesses-code-adventures/slop/db"
	_ "github.com/jesses-code-adventures/slop/env"
	"github.com/jesses-code-adventures/slop/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx := context.WithValue(context.Background(), "logger", logger)
	server := server.NewServer(ctx)
	db := db.NewPgDatabase(ctx)
	defer db.Close()
	api := a.NewAppApi(&server, ctx, db)
	a.BindRoutes(api)
	server.Serve()
}
