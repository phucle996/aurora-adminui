package main

import (
	"aurora-adminui/internal/app"
	"aurora-adminui/internal/config"
	"aurora-adminui/pkg/logger"
	"context"
	"log"
)

func main() {
	cfg := config.Load()
	logger.InitLogger(cfg)

	application, err := app.New(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
