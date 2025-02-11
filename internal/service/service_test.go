package service

import (
	"context"
	"ozon-tesk-task/internal/repository"
	"ozon-tesk-task/internal/service/mocks"
	"ozon-tesk-task/internal/transport/graph/model"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestService_ListPosts(t *testing.T) {
	type (
		mockBehavior func(r *mocks.Repository, limit, offset int32)
		args         struct {
			ctx          context.Context
			limit        int32
			offset       int32
			withComments bool
		}
	)

	posts := []*model.Post{
		{
			ID: 1,
		},
		{
			ID: 2,
		},
	}

	tests := []struct {
		name     string
		args     args
		repoMock mockBehavior
		want     []*model.Post
		wantErr  bool
	}{
		{
			name: "with comments test",
			args: args{
				ctx:          context.Background(),
				limit:        10,
				offset:       0,
				withComments: true,
			},
			repoMock: func(r *mocks.Repository, limit, offset int32) {
				r.On("ListPostsWithComments", mock.Anything, limit, offset).Return(posts, nil)
			},
			want:    posts,
			wantErr: false,
		},
		{
			name: "without comments test",
			args: args{
				ctx:          context.Background(),
				limit:        10,
				offset:       0,
				withComments: false,
			},
			repoMock: func(r *mocks.Repository, limit, offset int32) {
				r.On("ListPosts", mock.Anything, limit, offset).Return(posts, nil)
			},
			want:    posts,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			s := &Service{
				repo: r,
			}

			tt.repoMock(r, tt.args.limit, tt.args.offset)

			got, err := s.ListPosts(tt.args.ctx, tt.args.limit, tt.args.offset, tt.args.withComments)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ListPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ListPosts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetPostById(t *testing.T) {
	type (
		mockBehavior func(r *mocks.Repository, id int32)
		args         struct {
			ctx          context.Context
			id           int32
			withComments bool
		}
	)

	post := &model.Post{
		ID: 1,
	}

	tests := []struct {
		name     string
		args     args
		repoMock mockBehavior
		want     *model.Post
		wantErr  bool
	}{
		{
			name: "with comments test",
			args: args{
				ctx:          context.Background(),
				id:           1,
				withComments: true,
			},
			repoMock: func(r *mocks.Repository, id int32) {
				r.On("GetPostByIdWithComments", mock.Anything, id).Return(post, nil)
			},
			want:    post,
			wantErr: false,
		},
		{
			name: "without comments test",
			args: args{
				ctx:          context.Background(),
				id:           1,
				withComments: false,
			},
			repoMock: func(r *mocks.Repository, id int32) {
				r.On("GetPostById", mock.Anything, id).Return(post, nil)
			},
			want:    post,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			s := &Service{
				repo: r,
			}

			tt.repoMock(r, tt.args.id)

			got, err := s.GetPostById(tt.args.ctx, tt.args.id, tt.args.withComments)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetPostById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetPostById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_CreateComment(t *testing.T) {

	type (
		mockBehavior func(r *mocks.Repository, comment *model.Comment)
		args         struct {
			ctx     context.Context
			comment *model.Comment
		}
	)

	tests := []struct {
		name     string
		args     args
		repoMock mockBehavior
		want     *model.Comment
		wantErr  bool
	}{
		{
			name: "Wrong post id",
			args: args{
				ctx: context.Background(),
				comment: &model.Comment{
					PostID: 23131,
				},
			},
			repoMock: func(r *mocks.Repository, comment *model.Comment) {
				r.On("GetPostById", mock.Anything, comment.PostID).Return(nil, repository.ErrWrongPostId)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Comments not allowed",
			args: args{
				ctx: context.Background(),
				comment: &model.Comment{
					PostID: 1,
				},
			},
			repoMock: func(r *mocks.Repository, comment *model.Comment) {
				r.On("GetPostById", mock.Anything, comment.PostID).Return(&model.Post{AllowComments: false}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Parent does not exist",
			args: args{
				ctx: context.Background(),
				comment: &model.Comment{
					PostID:   1,
					ParentID: func() *int32 { v := int32(2132112); return &v }(),
				},
			},
			repoMock: func(r *mocks.Repository, comment *model.Comment) {
				r.On("GetPostById", mock.Anything, comment.PostID).Return(&model.Post{AllowComments: true}, nil)
				r.On("GetCommentById", mock.Anything, *comment.ParentID).Return(nil, repository.ErrWrongCommentId)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Parent is not in post",
			args: args{
				ctx: context.Background(),
				comment: &model.Comment{
					PostID:   1,
					ParentID: func() *int32 { v := int32(5); return &v }(),
				},
			},
			repoMock: func(r *mocks.Repository, comment *model.Comment) {
				r.On("GetPostById", mock.Anything, comment.PostID).Return(&model.Post{AllowComments: true}, nil)
				r.On("GetCommentById", mock.Anything, *comment.ParentID).Return(&model.Comment{PostID: 3}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "OK test",
			args: args{
				ctx: context.Background(),
				comment: &model.Comment{
					PostID:   1,
					ParentID: func() *int32 { v := int32(1); return &v }(),
				},
			},
			repoMock: func(r *mocks.Repository, comment *model.Comment) {
				r.On("GetPostById", mock.Anything, comment.PostID).Return(&model.Post{AllowComments: true}, nil)
				r.On("GetCommentById", mock.Anything, *comment.ParentID).Return(&model.Comment{PostID: 1}, nil)
				r.On("CreateComment", mock.Anything, comment).Return(int32(8), nil)
			},
			want: &model.Comment{
				ID:       8,
				PostID:   1,
				ParentID: func() *int32 { v := int32(1); return &v }(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			s := &Service{
				repo: r,
			}

			tt.repoMock(r, tt.args.comment)

			got, err := s.CreateComment(tt.args.ctx, tt.args.comment)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.CreateComment() = %v, want %v", got, tt.want)
			}
		})
	}
}
