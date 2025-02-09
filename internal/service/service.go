package service

import (
	"context"
	"ozon-tesk-task/internal/transport/graph/model"
)

type Repository interface {
	ListPosts(ctx context.Context, limit, offset int32) ([]*model.Post, error)
	ListPostsWithComments(ctx context.Context, limit, offset int32) ([]*model.Post, error)
	CreatePost(ctx context.Context, post *model.Post) (string, error)
	GetPostById(ctx context.Context, id string) (*model.Post, error)
	GetPostByIdWithComments(ctx context.Context, id string) (*model.Post, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) ListPosts(ctx context.Context, limit, offset int32, withComments bool) ([]*model.Post, error) {
	if withComments {
		return s.repo.ListPostsWithComments(ctx, limit, offset)
	}

	return s.repo.ListPosts(ctx, limit, offset)
}

func (s *Service) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	id, err := s.repo.CreatePost(ctx, post)
	if err != nil {
		return nil, err
	}

	post.ID = id

	return post, nil
}

func (s *Service) GetPostById(ctx context.Context, id string, withComments bool) (*model.Post, error) {
	if withComments {
		return s.repo.GetPostByIdWithComments(ctx, id)
	}
	return s.repo.GetPostById(ctx, id)
}
