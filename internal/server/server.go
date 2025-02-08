package server

import (
	"context"
	"fmt"
	"net/http"
	"ozon-tesk-task/internal/config"
	"ozon-tesk-task/pkg/logger"
	"time"
)

const shutdownTimeout = 5 * time.Second

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.ServicePort,
			Handler: handler,
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
