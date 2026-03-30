package main

import (
	"aurora-adminui/internal/app"
	"aurora-adminui/internal/config"
	"context"
	"log"
)

func main() {
	cfg := config.Load()

	application, err := app.New(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
