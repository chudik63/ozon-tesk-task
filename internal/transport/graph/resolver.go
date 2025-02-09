package graph

import (
	"context"
	"ozon-tesk-task/internal/transport/graph/model"
	"ozon-tesk-task/pkg/logger"
)

type Service interface {
	ListPosts(ctx context.Context, limit, offset int32) ([]*model.Post, error)
	// GetPost(id string) *model.Post
}

type Resolver struct {
	service Service
	logs    logger.Logger
}

func NewResolver(srv Service, logs logger.Logger) *Resolver {
	return &Resolver{
		service: srv,
		logs:    logs,
	}
}
