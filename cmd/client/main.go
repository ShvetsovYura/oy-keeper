package main

import (
	"context"

	"github.com/ShvetsovYura/oykeeper/internal/client"
	"github.com/ShvetsovYura/oykeeper/internal/logger"
)

func main() {
	logger.Init()
	c := client.New(":9091")
	logger.Log.Info("startging up client app")
	c.CreateRecord(context.Background())
}
