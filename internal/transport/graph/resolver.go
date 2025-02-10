package graph

import (
	"context"
	"ozon-tesk-task/internal/transport/graph/model"
	"ozon-tesk-task/pkg/logger"
)

type Service interface {
	ListPosts(ctx context.Context, limit, offset int32, withComments bool) ([]*model.Post, error)
	CreatePost(ctx context.Context, post *model.Post) (*model.Post, error)
	GetPostById(ctx context.Context, id int32, withComments bool) (*model.Post, error)
	CreateComment(ctx context.Context, comment *model.Comment) (*model.Comment, error)
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
