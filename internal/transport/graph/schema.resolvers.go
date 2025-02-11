package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.64

import (
	"context"
	"errors"
	"net/http"
	"ozon-tesk-task/internal/preloads"
	"ozon-tesk-task/internal/repository"
	"ozon-tesk-task/internal/transport/graph/model"
	"ozon-tesk-task/pkg/pointer"
	"time"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
	if input.Title == "" || input.Content == "" {
		r.logs.Info(ctx, "invalid input arguments")
		return nil, &gqlerror.Error{
			Message: "invalid argument",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	if len(input.Title) > 200 {
		r.logs.Info(ctx, "input title is too long")
		return nil, &gqlerror.Error{
			Message: "title is too long",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	if len(input.Content) > 2000 {
		r.logs.Info(ctx, "input content is too long")
		return nil, &gqlerror.Error{
			Message: "content is too long",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	r.logs.Debug(ctx, "Creating post", zap.Any("input", input))

	post, err := r.service.CreatePost(ctx, &model.Post{
		Title:         input.Title,
		Content:       input.Content,
		AllowComments: input.AllowComments,
		CreatedAt:     time.Now().Format(time.DateTime),
		Author:        0,
	})
	if err != nil {
		r.logs.Error(ctx, "failed to create post", zap.String("err", err.Error()))
		return nil, &gqlerror.Error{
			Message: "failed to create post",
			Extensions: map[string]interface{}{
				"code": http.StatusInternalServerError,
			},
		}
	}

	return post, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.CreateCommentInput) (*model.Comment, error) {
	if len(input.Content) > 2000 {
		r.logs.Info(ctx, "comment is too long")
		return nil, &gqlerror.Error{
			Message: "comment is too long",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	if input.Content == "" || input.PostID <= 0 || pointer.Deref(input.ParentID, 0) < 0 {
		r.logs.Info(ctx, "invalid input arguments")
		return nil, &gqlerror.Error{
			Message: "invalid argument",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	r.logs.Debug(ctx, "Creating comment", zap.Any("input", input))

	comment, err := r.service.CreateComment(ctx, &model.Comment{
		PostID:    input.PostID,
		ParentID:  input.ParentID,
		Content:   input.Content,
		CreatedAt: time.Now().Format(time.DateTime),
		Author:    0,
	})

	if err != nil {
		if errors.Is(err, repository.ErrCommentsNotAllowed) {
			r.logs.Error(ctx, "can`t create comment", zap.String("err", err.Error()))
			return nil, &gqlerror.Error{
				Message: err.Error(),
				Extensions: map[string]interface{}{
					"code": http.StatusForbidden,
				},
			}
		}
		if errors.Is(err, repository.ErrWrongCommentId) {
			r.logs.Error(ctx, "can`t create comment", zap.String("err", err.Error()))
			return nil, &gqlerror.Error{
				Message: err.Error(),
				Extensions: map[string]interface{}{
					"code": http.StatusBadRequest,
				},
			}
		}
		if errors.Is(err, repository.ErrMatchCommentWithPost) {
			r.logs.Error(ctx, "can`t create comment", zap.String("err", err.Error()))
			return nil, &gqlerror.Error{
				Message: err.Error(),
				Extensions: map[string]interface{}{
					"code": http.StatusBadRequest,
				},
			}
		}

		r.logs.Error(ctx, "failed to create comment", zap.String("err", err.Error()))
		return nil, &gqlerror.Error{
			Message: "failed to create comment",
			Extensions: map[string]interface{}{
				"code": http.StatusInternalServerError,
			},
		}
	}

	r.pubsub.Publish(ctx, comment)

	return comment, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, page *int32, limit *int32) ([]*model.Post, error) {
	var (
		withComments bool
	)

	pl := preloads.GetPreloads(ctx)

	for _, field := range pl {
		if field == "comments" {
			withComments = true
			break
		}
	}

	lim := pointer.Deref(limit, 10)
	p := pointer.Deref(page, 1)

	if lim <= 0 {
		r.logs.Warn(ctx, "invalid pagination argument", zap.Int32("limit", lim))
		lim = 10
	}
	if p < 1 {
		r.logs.Warn(ctx, "invalid pagination argument", zap.Int32("page", p))
		p = 1
	}

	offset := lim * (p - 1)

	r.logs.Debug(ctx, "Loading posts", zap.Bool("with comments", withComments), zap.Int32("page", p))

	posts, err := r.service.ListPosts(ctx, lim, offset, withComments)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			r.logs.Error(ctx, "can`t list posts", zap.String("err", err.Error()))
			return nil, &gqlerror.Error{
				Message: err.Error(),
				Extensions: map[string]interface{}{
					"code": http.StatusNotFound,
				},
			}
		}

		r.logs.Error(ctx, "failed to list posts", zap.String("err", err.Error()))
		return nil, &gqlerror.Error{
			Message: "failed to list posts",
			Extensions: map[string]interface{}{
				"code": http.StatusInternalServerError,
			},
		}
	}

	return posts, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id int32) (*model.Post, error) {
	if id <= 0 {
		return nil, &gqlerror.Error{
			Message: "invalid argument",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	var (
		withComments bool
	)

	pl := preloads.GetPreloads(ctx)

	for _, field := range pl {
		if field == "comments" {
			withComments = true
			break
		}
	}

	r.logs.Debug(ctx, "Loading post", zap.Int32("id", id), zap.Bool("with comments", withComments))

	post, err := r.service.GetPostById(ctx, id, withComments)
	if err != nil {
		if errors.Is(err, repository.ErrWrongPostId) {
			r.logs.Error(ctx, "can`t get post", zap.String("err", err.Error()))
			return nil, &gqlerror.Error{
				Message: err.Error(),
				Extensions: map[string]interface{}{
					"code": http.StatusNotFound,
				},
			}
		}

		r.logs.Error(ctx, "failed to get post", zap.String("err", err.Error()))
		return nil, &gqlerror.Error{
			Message: "failed to get post",
			Extensions: map[string]interface{}{
				"code": http.StatusInternalServerError,
			},
		}
	}

	return post, nil
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID int32, page *int32, limit *int32) ([]*model.Comment, error) {

	lim := pointer.Deref(limit, 10)
	p := pointer.Deref(page, 1)

	if lim <= 0 {
		r.logs.Warn(ctx, "invalid pagination argument", zap.Int32("limit", lim))
		lim = 10
	}
	if p < 1 {
		r.logs.Warn(ctx, "invalid pagination argument", zap.Int32("page", p))
		p = 1
	}

	offset := lim * (p - 1)

	r.logs.Debug(ctx, "Loading comments", zap.Int32("post", postID), zap.Int32("page", p))

	comments, err := r.service.GetComments(ctx, postID, lim, offset)
	if err != nil {
		if errors.Is(err, repository.ErrWrongPostId) {
			r.logs.Error(ctx, "can`t get post", zap.String("err", err.Error()))
			return nil, &gqlerror.Error{
				Message: err.Error(),
				Extensions: map[string]interface{}{
					"code": http.StatusNotFound,
				},
			}
		}

		r.logs.Error(ctx, "failed to get post", zap.String("err", err.Error()))
		return nil, &gqlerror.Error{
			Message: "failed to get post",
			Extensions: map[string]interface{}{
				"code": http.StatusInternalServerError,
			},
		}
	}

	return comments, nil
}

// DeletePost is the resolver for the deletePost field.
func (r *queryResolver) DeletePost(ctx context.Context, postID int32) (int32, error) {
	if postID <= 0 {
		return 0, &gqlerror.Error{
			Message: "invalid argument",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	r.logs.Debug(ctx, "Deleting post", zap.Int32("id", postID))

	err := r.service.DeletePost(ctx, postID)
	if err != nil {
		if errors.Is(err, repository.ErrWrongPostId) {
			r.logs.Error(ctx, "can`t get post", zap.String("err", err.Error()))
			return 0, &gqlerror.Error{
				Message: err.Error(),
				Extensions: map[string]interface{}{
					"code": http.StatusNotFound,
				},
			}
		}

		r.logs.Error(ctx, "failed to delete post", zap.String("err", err.Error()))
		return 0, &gqlerror.Error{
			Message: "failed to delete post",
			Extensions: map[string]interface{}{
				"code": http.StatusInternalServerError,
			},
		}
	}

	return postID, nil
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID int32) (<-chan *model.Comment, error) {
	if postID <= 0 {
		return nil, &gqlerror.Error{
			Message: "invalid argument",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	if !r.pubsub.Check(postID) {
		_, err := r.service.GetPostById(ctx, postID, false)
		if err != nil {
			return nil, err
		}
	}

	r.logs.Debug(ctx, "Creating new subscription", zap.Int32("postId", postID))

	ch := r.pubsub.Subscribe(ctx, postID)

	return ch, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
