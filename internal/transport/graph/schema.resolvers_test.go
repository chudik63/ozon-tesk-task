package graph

import (
	"context"
	"errors"
	"ozon-tesk-task/internal/repository"
	"ozon-tesk-task/internal/transport/graph/mocks"
	"ozon-tesk-task/internal/transport/graph/model"
	"ozon-tesk-task/pkg/logger"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func Test_mutationResolver_CreatePost(t *testing.T) {
	type (
		mockServiceBehavior func(s *mocks.Service, post, returnPost *model.Post)
		args                struct {
			ctx   context.Context
			input model.CreatePostInput
		}
	)

	tests := []struct {
		name        string
		args        args
		want        *model.Post
		serviceMock mockServiceBehavior
		wantErr     bool
	}{
		{
			name: "OK test",
			args: args{
				ctx: context.Background(),
				input: model.CreatePostInput{
					Title:         "title",
					Content:       "content",
					AllowComments: true,
				},
			},
			want: &model.Post{
				ID:            1,
				Title:         "title",
				Content:       "content",
				AllowComments: true,
			},
			serviceMock: func(s *mocks.Service, post *model.Post, returnPost *model.Post) {
				s.On("CreatePost", mock.Anything, post).Return(returnPost, nil)
			},
			wantErr: false,
		},
		{
			name: "Empty title",
			args: args{
				ctx: context.Background(),
				input: model.CreatePostInput{
					Title:         "",
					Content:       "content",
					AllowComments: true,
				},
			},
			want:        nil,
			serviceMock: func(s *mocks.Service, post *model.Post, returnPost *model.Post) {},
			wantErr:     true,
		},
		{
			name: "Empty сontent",
			args: args{
				ctx: context.Background(),
				input: model.CreatePostInput{
					Title:         "title",
					Content:       "",
					AllowComments: true,
				},
			},
			want:        nil,
			serviceMock: func(s *mocks.Service, post *model.Post, returnPost *model.Post) {},
			wantErr:     true,
		},
		{
			name: "Internal Error",
			args: args{
				ctx: context.Background(),
				input: model.CreatePostInput{
					Title:         "title",
					Content:       "content",
					AllowComments: true,
				},
			},
			want: nil,
			serviceMock: func(s *mocks.Service, post *model.Post, returnPost *model.Post) {
				s.On("CreatePost", mock.Anything, post).Return(nil, errors.New("internal error"))
			},
			wantErr: true,
		},
		{
			name: "Long content",
			args: args{
				ctx: context.Background(),
				input: model.CreatePostInput{
					Title:         "title",
					Content:       string(make([]byte, 2001)),
					AllowComments: true,
				},
			},
			want:        nil,
			serviceMock: func(s *mocks.Service, post *model.Post, returnPost *model.Post) {},
			wantErr:     true,
		},
		{
			name: "Long title",
			args: args{
				ctx: context.Background(),
				input: model.CreatePostInput{
					Title:         string(make([]byte, 201)),
					Content:       "content",
					AllowComments: true,
				},
			},
			want:        nil,
			serviceMock: func(s *mocks.Service, post *model.Post, returnPost *model.Post) {},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mocks.NewService(t)
			log, _ := logger.New("test")
			p := mocks.NewPubSub(t)

			r := &mutationResolver{
				Resolver: &Resolver{s, log, p},
			}

			tt.serviceMock(s, &model.Post{
				Title:         tt.args.input.Title,
				Content:       tt.args.input.Content,
				AllowComments: tt.args.input.AllowComments,
				CreatedAt:     time.Now().Format(time.DateTime),
				Author:        0,
			}, tt.want)

			got, err := r.CreatePost(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("mutationResolver.CreatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mutationResolver.CreatePost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mutationResolver_CreateComment(t *testing.T) {
	type (
		mockServiceBehavior func(s *mocks.Service, comment, returnComment *model.Comment)
		pubSubBehavior      func(p *mocks.PubSub, comment *model.Comment)
		args                struct {
			ctx   context.Context
			input model.CreateCommentInput
		}
	)

	tests := []struct {
		name        string
		args        args
		mockService mockServiceBehavior
		mockPubSub  pubSubBehavior
		want        *model.Comment
		wantErr     bool
	}{
		{
			name: "OK test",
			args: args{
				ctx: context.Background(),
				input: model.CreateCommentInput{
					PostID:  1,
					Content: "content",
				},
			},
			mockService: func(s *mocks.Service, comment, returnComment *model.Comment) {
				s.On("CreateComment", mock.Anything, comment).Return(returnComment, nil)
			},
			mockPubSub: func(p *mocks.PubSub, comment *model.Comment) {
				p.On("Publish", mock.Anything, comment)
			},
			want: &model.Comment{
				ID:      1,
				PostID:  1,
				Content: "content",
			},
			wantErr: false,
		},
		{
			name: "With parent id",
			args: args{
				ctx: context.Background(),
				input: model.CreateCommentInput{
					PostID:   1,
					ParentID: func() *int32 { v := int32(1); return &v }(),
					Content:  "content",
				},
			},
			mockService: func(s *mocks.Service, comment, returnComment *model.Comment) {
				s.On("CreateComment", mock.Anything, comment).Return(returnComment, nil)
			},
			mockPubSub: func(p *mocks.PubSub, comment *model.Comment) {
				p.On("Publish", mock.Anything, comment)
			},
			want: &model.Comment{
				ID:      1,
				PostID:  1,
				Content: "content",
			},
			wantErr: false,
		},
		{
			name: "Invalid post id",
			args: args{
				ctx: context.Background(),
				input: model.CreateCommentInput{
					PostID:  -1,
					Content: "content",
				},
			},
			mockService: func(s *mocks.Service, comment, returnComment *model.Comment) {},
			mockPubSub:  func(p *mocks.PubSub, comment *model.Comment) {},
			want:        nil,
			wantErr:     true,
		},
		{
			name: "Long content",
			args: args{
				ctx: context.Background(),
				input: model.CreateCommentInput{
					PostID:  1,
					Content: string(make([]byte, 2001)),
				},
			},
			mockService: func(s *mocks.Service, comment, returnComment *model.Comment) {},
			mockPubSub:  func(p *mocks.PubSub, comment *model.Comment) {},
			want:        nil,
			wantErr:     true,
		},
		{
			name: "Invalid post id",
			args: args{
				ctx: context.Background(),
				input: model.CreateCommentInput{
					PostID:   1,
					ParentID: func() *int32 { v := int32(-1); return &v }(),
					Content:  "content",
				},
			},
			mockService: func(s *mocks.Service, comment, returnComment *model.Comment) {},
			mockPubSub:  func(p *mocks.PubSub, comment *model.Comment) {},
			want:        nil,
			wantErr:     true,
		},
		{
			name: "Wrong post id",
			args: args{
				ctx: context.Background(),
				input: model.CreateCommentInput{
					PostID:  1232131212,
					Content: "content",
				},
			},
			mockService: func(s *mocks.Service, comment, returnComment *model.Comment) {
				s.On("CreateComment", mock.Anything, comment).Return(nil, repository.ErrWrongPostId)
			},
			mockPubSub: func(p *mocks.PubSub, comment *model.Comment) {},
			want:       nil,
			wantErr:    true,
		},
		{
			name: "Wrong parent id",
			args: args{
				ctx: context.Background(),
				input: model.CreateCommentInput{
					PostID:   1,
					ParentID: func() *int32 { v := int32(1233231); return &v }(),
					Content:  "content",
				},
			},
			mockService: func(s *mocks.Service, comment, returnComment *model.Comment) {
				s.On("CreateComment", mock.Anything, comment).Return(nil, repository.ErrWrongCommentId)
			},
			mockPubSub: func(p *mocks.PubSub, comment *model.Comment) {},
			want:       nil,
			wantErr:    true,
		},
		{
			name: "Parent is not in Post",
			args: args{
				ctx: context.Background(),
				input: model.CreateCommentInput{
					PostID:   1,
					ParentID: func() *int32 { v := int32(5); return &v }(),
					Content:  "content",
				},
			},
			mockService: func(s *mocks.Service, comment, returnComment *model.Comment) {
				s.On("CreateComment", mock.Anything, comment).Return(nil, repository.ErrMatchCommentWithPost)
			},
			mockPubSub: func(p *mocks.PubSub, comment *model.Comment) {},
			want:       nil,
			wantErr:    true,
		},
		{
			name: "Comments not allowed",
			args: args{
				ctx: context.Background(),
				input: model.CreateCommentInput{
					PostID:   1,
					ParentID: func() *int32 { v := int32(1); return &v }(),
					Content:  "content",
				},
			},
			mockService: func(s *mocks.Service, comment, returnComment *model.Comment) {
				s.On("CreateComment", mock.Anything, comment).Return(nil, repository.ErrCommentsNotAllowed)
			},
			mockPubSub: func(p *mocks.PubSub, comment *model.Comment) {},
			want:       nil,
			wantErr:    true,
		},
		{
			name: "Internal error",
			args: args{
				ctx: context.Background(),
				input: model.CreateCommentInput{
					PostID:   1,
					ParentID: func() *int32 { v := int32(1); return &v }(),
					Content:  "content",
				},
			},
			mockService: func(s *mocks.Service, comment, returnComment *model.Comment) {
				s.On("CreateComment", mock.Anything, comment).Return(nil, errors.New("internal error"))
			},
			mockPubSub: func(p *mocks.PubSub, comment *model.Comment) {},
			want:       nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mocks.NewService(t)
			log, _ := logger.New("test")
			p := mocks.NewPubSub(t)

			r := &mutationResolver{
				Resolver: &Resolver{s, log, p},
			}

			tt.mockService(s, &model.Comment{
				PostID:    tt.args.input.PostID,
				ParentID:  tt.args.input.ParentID,
				Content:   tt.args.input.Content,
				CreatedAt: time.Now().Format(time.DateTime),
				Author:    0,
			}, tt.want)

			tt.mockPubSub(p, tt.want)

			got, err := r.CreateComment(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("mutationResolver.CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mutationResolver.CreateComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_queryResolver_Posts(t *testing.T) {
	type (
		mockServiceBehavior func(s *mocks.Service, returnPosts []*model.Post)
		args                struct {
			ctx   context.Context
			page  *int32
			limit *int32
		}
	)

	tests := []struct {
		name        string
		args        args
		serviceMock mockServiceBehavior
		want        []*model.Post
		wantErr     bool
	}{
		{
			name: "OK test",
			args: args{
				ctx:   context.Background(),
				page:  func() *int32 { v := int32(1); return &v }(),
				limit: func() *int32 { v := int32(10); return &v }(),
			},
			serviceMock: func(s *mocks.Service, returnPosts []*model.Post) {
				s.On("ListPosts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(returnPosts, nil)
			},
			want: []*model.Post{
				{
					ID:    1,
					Title: "empty",
				},
				{
					ID:    2,
					Title: "empty",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid page value",
			args: args{
				ctx:   context.Background(),
				page:  func() *int32 { v := int32(-1); return &v }(),
				limit: func() *int32 { v := int32(10); return &v }(),
			},
			serviceMock: func(s *mocks.Service, returnPosts []*model.Post) {
				s.On("ListPosts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(returnPosts, nil)
			},
			want: []*model.Post{
				{
					ID:    1,
					Title: "empty",
				},
				{
					ID:    2,
					Title: "empty",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid limit value",
			args: args{
				ctx:   context.Background(),
				page:  func() *int32 { v := int32(1); return &v }(),
				limit: func() *int32 { v := int32(0); return &v }(),
			},
			serviceMock: func(s *mocks.Service, returnPosts []*model.Post) {
				s.On("ListPosts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(returnPosts, nil)
			},
			want: []*model.Post{
				{
					ID:    1,
					Title: "empty",
				},
				{
					ID:    2,
					Title: "empty",
				},
			},
			wantErr: false,
		},
		{
			name: "List first page of 1 element",
			args: args{
				ctx:   context.Background(),
				page:  func() *int32 { v := int32(1); return &v }(),
				limit: func() *int32 { v := int32(1); return &v }(),
			},
			serviceMock: func(s *mocks.Service, returnPosts []*model.Post) {
				s.On("ListPosts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(returnPosts, nil)
			},
			want: []*model.Post{
				{
					ID:    1,
					Title: "empty",
				},
			},
			wantErr: false,
		},
		{
			name: "List second page of 1 element",
			args: args{
				ctx:   context.Background(),
				page:  func() *int32 { v := int32(2); return &v }(),
				limit: func() *int32 { v := int32(1); return &v }(),
			},
			serviceMock: func(s *mocks.Service, returnPosts []*model.Post) {
				s.On("ListPosts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(returnPosts, nil)
			},
			want: []*model.Post{
				{
					ID:    1,
					Title: "empty",
				},
			},
			wantErr: false,
		},
		{
			name: "List unexisted page",
			args: args{
				ctx:   context.Background(),
				page:  func() *int32 { v := int32(3); return &v }(),
				limit: func() *int32 { v := int32(1); return &v }(),
			},
			serviceMock: func(s *mocks.Service, returnPosts []*model.Post) {
				s.On("ListPosts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, repository.ErrNotFound)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Internal error",
			args: args{
				ctx:   context.Background(),
				page:  func() *int32 { v := int32(3); return &v }(),
				limit: func() *int32 { v := int32(1); return &v }(),
			},
			serviceMock: func(s *mocks.Service, returnPosts []*model.Post) {
				s.On("ListPosts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("internal error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mocks.NewService(t)
			log, _ := logger.New("test")
			p := mocks.NewPubSub(t)

			r := &queryResolver{
				Resolver: &Resolver{s, log, p},
			}

			tt.serviceMock(s, tt.want)

			got, err := r.Posts(tt.args.ctx, tt.args.page, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("queryResolver.Posts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryResolver.Posts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_queryResolver_Post(t *testing.T) {
	type (
		mockServiceBehavior func(s *mocks.Service, id int32, returnPost *model.Post)
		args                struct {
			ctx context.Context
			id  int32
		}
	)

	tests := []struct {
		name        string
		args        args
		serviceMock mockServiceBehavior
		want        *model.Post
		wantErr     bool
	}{
		{
			name: "OK test",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			serviceMock: func(s *mocks.Service, id int32, returnPost *model.Post) {
				s.On("GetPostById", mock.Anything, id, mock.Anything).Return(returnPost, nil)
			},
			want: &model.Post{
				ID:    1,
				Title: "test",
			},
			wantErr: false,
		},
		{
			name: "Invalid post id",
			args: args{
				ctx: context.Background(),
				id:  -1,
			},
			serviceMock: func(s *mocks.Service, id int32, returnPost *model.Post) {},
			want:        nil,
			wantErr:     true,
		},
		{
			name: "Wrong post id",
			args: args{
				ctx: context.Background(),
				id:  213123213,
			},
			serviceMock: func(s *mocks.Service, id int32, returnPost *model.Post) {
				s.On("GetPostById", mock.Anything, id, mock.Anything).Return(nil, repository.ErrWrongPostId)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Internal error",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			serviceMock: func(s *mocks.Service, id int32, returnPost *model.Post) {
				s.On("GetPostById", mock.Anything, id, mock.Anything).Return(nil, errors.New("internal error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mocks.NewService(t)
			log, _ := logger.New("test")
			p := mocks.NewPubSub(t)

			r := &queryResolver{
				Resolver: &Resolver{s, log, p},
			}

			tt.serviceMock(s, tt.args.id, tt.want)

			got, err := r.Post(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("queryResolver.Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryResolver.Post() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_queryResolver_DeletePost(t *testing.T) {
	type (
		mockServiceBehavior func(s *mocks.Service, postID int32)
		args                struct {
			ctx    context.Context
			postID int32
		}
	)

	tests := []struct {
		name        string
		args        args
		serviceMock mockServiceBehavior
		want        int32
		wantErr     bool
	}{
		{
			name: "OK test",
			args: args{
				ctx:    context.Background(),
				postID: 1,
			},
			serviceMock: func(s *mocks.Service, postID int32) {
				s.On("DeletePost", mock.Anything, postID).Return(nil)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Invalid post id",
			args: args{
				ctx:    context.Background(),
				postID: -1,
			},
			serviceMock: func(s *mocks.Service, postID int32) {},
			want:        0,
			wantErr:     true,
		},
		{
			name: "Wrong post id",
			args: args{
				ctx:    context.Background(),
				postID: 3123213,
			},
			serviceMock: func(s *mocks.Service, postID int32) {
				s.On("DeletePost", mock.Anything, postID).Return(repository.ErrWrongPostId)
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Internal error",
			args: args{
				ctx:    context.Background(),
				postID: 1,
			},
			serviceMock: func(s *mocks.Service, postID int32) {
				s.On("DeletePost", mock.Anything, postID).Return(errors.New("internal error"))
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mocks.NewService(t)
			log, _ := logger.New("test")
			p := mocks.NewPubSub(t)

			r := &queryResolver{
				Resolver: &Resolver{s, log, p},
			}

			tt.serviceMock(s, tt.args.postID)

			got, err := r.DeletePost(tt.args.ctx, tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("queryResolver.DeletePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("queryResolver.DeletePost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_subscriptionResolver_CommentAdded(t *testing.T) {
	type (
		mockServiceBehavior func(s *mocks.Service, postID int32)
		pubSubBehavior      func(p *mocks.PubSub, postId int32)
		args                struct {
			ctx    context.Context
			postID int32
		}
	)

	ch := make(<-chan *model.Comment, 1)

	tests := []struct {
		name        string
		args        args
		serviceMock mockServiceBehavior
		pubsubMock  pubSubBehavior
		want        <-chan *model.Comment
		wantErr     bool
	}{
		{
			name: "Invalid post id",
			args: args{
				ctx:    context.Background(),
				postID: -1,
			},
			pubsubMock:  func(p *mocks.PubSub, postId int32) {},
			serviceMock: func(s *mocks.Service, postID int32) {},
			want:        nil,
			wantErr:     true,
		},
		{
			name: "Wrong post id",
			args: args{
				ctx:    context.Background(),
				postID: 312313132,
			},
			pubsubMock: func(p *mocks.PubSub, postId int32) {
				p.On("Check", postId).Return(false)
			},
			serviceMock: func(s *mocks.Service, postID int32) {
				s.On("GetPostById", mock.Anything, postID, mock.Anything).Return(nil, repository.ErrWrongPostId)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Post does not have subscribers",
			args: args{
				ctx:    context.Background(),
				postID: 1,
			},
			pubsubMock: func(p *mocks.PubSub, postId int32) {
				p.On("Check", postId).Return(false)
				p.On("Subscribe", mock.Anything, postId).Return(ch)
			},
			serviceMock: func(s *mocks.Service, postID int32) {
				s.On("GetPostById", mock.Anything, postID, mock.Anything).Return(&model.Post{}, nil)
			},
			want:    ch,
			wantErr: false,
		},
		{
			name: "Post already has subscribers",
			args: args{
				ctx:    context.Background(),
				postID: 1,
			},
			pubsubMock: func(p *mocks.PubSub, postId int32) {
				p.On("Check", postId).Return(true)
				p.On("Subscribe", mock.Anything, postId).Return(ch)
			},
			serviceMock: func(s *mocks.Service, postID int32) {},
			want:        ch,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mocks.NewService(t)
			log, _ := logger.New("test")
			p := mocks.NewPubSub(t)

			r := &subscriptionResolver{
				Resolver: &Resolver{s, log, p},
			}

			tt.serviceMock(s, tt.args.postID)
			tt.pubsubMock(p, tt.args.postID)

			got, err := r.CommentAdded(tt.args.ctx, tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("subscriptionResolver.CommentAdded() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subscriptionResolver.CommentAdded() = %v, want %v", got, tt.want)
			}
		})
	}
}
