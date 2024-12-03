package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/ShvetsovYura/oykeeper/internal/logger"
	"github.com/ShvetsovYura/oykeeper/internal/server"
)

func main() {
	logger.Init()
	s := server.New()
	logger.Log.Info("starting grpc server")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	s.Run(ctx)
}
