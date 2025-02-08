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
	"ozon-tesk-task/internal/transport/graph"
	"ozon-tesk-task/pkg/logger"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/vektah/gqlparser/v2/ast"
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

	router := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(service)}))
	router.AddTransport(transport.Options{})
	router.AddTransport(transport.GET{})
	router.AddTransport(transport.POST{})

	router.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	router.Use(extension.Introspection{})
	router.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	srv := server.NewServer(cfg, router)

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
