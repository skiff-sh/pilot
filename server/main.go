package main

import (
	"context"
	"log/slog"

	"github.com/skiff-sh/config"
	"github.com/skiff-sh/config/contexts"
	pilotconfig "github.com/skiff-sh/pilot/server/pkg/config"
	"github.com/skiff-sh/pilot/server/pkg/server"
)

func main() {
	conf, err := pilotconfig.New()
	if err != nil {
		panic(err)
	}

	logger, err := config.NewLogger(conf.Log)
	if err != nil {
		panic(err)
	}

	srv, err := server.New(conf, logger)
	if err != nil {
		logger.Error("Failed to create server.", slog.String("err", err.Error()))
		panic(err)
	}

	ctx := context.TODO()
	ctx = contexts.WithLogger(ctx, logger)
	logger.Info("GRPC started.", slog.String("addr", conf.GRPC.Addr.String()))
	err = srv.Start(ctx)
	if err != nil {
		logger.Error("Failed to start server.", slog.String("err", err.Error()))
		panic(err)
	}
}
