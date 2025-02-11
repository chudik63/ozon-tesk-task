package graph

import (
	"context"
	"ozon-tesk-task/internal/transport/graph/model"
	"ozon-tesk-task/pkg/logger"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name Service
type Service interface {
	ListPosts(ctx context.Context, limit, offset int32, withComments bool) ([]*model.Post, error)
	CreatePost(ctx context.Context, post *model.Post) (*model.Post, error)
	GetPostById(ctx context.Context, id int32, withComments bool) (*model.Post, error)
	GetComments(ctx context.Context, postId int32, limit, offset int32) ([]*model.Comment, error)
	CreateComment(ctx context.Context, comment *model.Comment) (*model.Comment, error)
	DeletePost(ctx context.Context, postId int32) error
}

//go:generate go run github.com/vektra/mockery/v2@latest --name PubSub
type PubSub interface {
	Subscribe(ctx context.Context, postId int32) <-chan *model.Comment
	Unsubscribe(ctx context.Context, postId int32, ch chan *model.Comment)
	Publish(ctx context.Context, comment *model.Comment)
	Check(postId int32) bool
}

type Resolver struct {
	service Service
	logs    logger.Logger
	pubsub  PubSub
}

func NewResolver(srv Service, logs logger.Logger, pubsub PubSub) *Resolver {
	return &Resolver{
		service: srv,
		logs:    logs,
		pubsub:  pubsub,
	}
}
