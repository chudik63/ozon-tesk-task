package main

import (
	"context"
	"fmt"
	"ozon-tesk-task/internal/config"
	"ozon-tesk-task/internal/database/postgres"
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

	db, err := postgres.New(cfg.Config)
	if err != nil {
		mainLogger.Fatal(ctx, err.Error())
	}

	err = postgres.Migrate(cfg.Config, cfg.MigrationsPath)
	if err != nil {
		mainLogger.Fatal(ctx, err.Error())
	}

	_ = db
}
