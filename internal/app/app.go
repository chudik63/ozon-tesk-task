package app

import (
	"context"
	"os"
	"os/signal"
	"ozon-tesk-task/internal/config"
	"ozon-tesk-task/internal/database"
	"ozon-tesk-task/internal/repository"
	"ozon-tesk-task/internal/server"
	"ozon-tesk-task/internal/service"
	"ozon-tesk-task/pkg/logger"
	"syscall"

	"go.uber.org/zap"
)

func Run(ctx context.Context, cfg *config.Config) {
	mainLogger := logger.GetLoggerFromCtx(ctx)

	mainLogger.Debug(ctx, "Storage type: "+cfg.StorageType)

	db, err := database.NewDatabase(ctx, cfg)
	if err != nil {
		mainLogger.Fatal(ctx, err.Error())
	}

	repo := repository.New(db)

	service := service.New(repo)

	_ = service

	srv := server.NewServer(cfg)

	go func() {
		if err := srv.Run(); err != nil {
			mainLogger.Fatal(ctx, "failed to run server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	signal := <-quit

	mainLogger.Debug(ctx, "Gracefully stopping the server", zap.String("caught signal", signal.String()))

	if err := srv.Stop(); err != nil {
		mainLogger.Error(ctx, "failed to stop server", zap.String("err", err.Error()))
	}
}
