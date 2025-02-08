package server

import (
	"context"
	"fmt"
	"net/http"
	"ozon-tesk-task/internal/config"
	"ozon-tesk-task/pkg/logger"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const shutdownTimeout = 5 * time.Second

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler *handler.Server) *Server {
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", handler)

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.ServicePort,
			Handler: nil,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	logs := logger.GetLoggerFromCtx(ctx)
	logs.Info(ctx, fmt.Sprintf("Starting server on %s", s.httpServer.Addr))

	err := s.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	ctx, shutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdown()

	return s.httpServer.Shutdown(ctx)
}
