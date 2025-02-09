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
	"ozon-tesk-task/internal/transport/http"
	"ozon-tesk-task/pkg/logger"
	"syscall"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func Run(ctx context.Context, cfg *config.Config) {
	mainLogger := logger.GetLoggerFromCtx(ctx)

	mainLogger.Debug(ctx, "Storage picked", zap.String("type", cfg.StorageType))

	db, err := database.NewDatabase(ctx, cfg)
	if err != nil {
		mainLogger.Fatal(ctx, err.Error())
	}

	repo := repository.New(db)

	service := service.New(repo)

	e := echo.New()

	http.NewHandler(e, service, mainLogger)

	srv := server.NewServer(cfg, e.Server.Handler)

	go func() {
		if err := srv.Run(ctx); err != nil {
			mainLogger.Fatal(ctx, "failed to run server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	signal := <-quit

	mainLogger.Info(ctx, "Gracefully stopping the server", zap.String("caught signal", signal.String()))

	if err := srv.Stop(); err != nil {
		mainLogger.Error(ctx, "failed to stop server", zap.String("err", err.Error()))
	}
}
