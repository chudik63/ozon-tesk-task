package main

import (
	"context"
	"fmt"
	"ozon-tesk-task/internal/app"
	"ozon-tesk-task/internal/config"
	"ozon-tesk-task/pkg/logger"
)

const (
	serviceName = "ozon-test-service"
)

func main() {
	mainLogger, err := logger.New(serviceName)
	if err != nil {
		panic(fmt.Errorf("failed to create logger: %v", err))
	}

	ctx := context.WithValue(context.Background(), logger.LoggerKey, mainLogger)

	cfg, err := config.New()
	if err != nil {
		mainLogger.Fatal(ctx, err.Error())
	}

	app.Run(ctx, cfg)
}
