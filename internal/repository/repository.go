package repository

import (
	"context"
	"database/sql"
	"errors"
	"ozon-tesk-task/internal/database"
	"ozon-tesk-task/internal/transport/graph/model"
	"strconv"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	db *database.Database
}

func New(db *database.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListPosts(ctx context.Context, limit, offset int32) ([]*model.Post, error) {
	rows, err := sq.Select("id", "user_id", "title", "content", "comments_allowed", "created_at").
		From("posts").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db.DB).
		Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post

	for rows.Next() {
		var post model.Post

		if err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.AllowComments, &post.CreatedAt); err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	if rows.Err() != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, ErrNotFound
	}

	return posts, nil
}

func (r *Repository) ListPostsWithComments(ctx context.Context, limit, offset int32) ([]*model.Post, error) {
	rows, err := sq.Select("p.id", "p.user_id", "p.title", "p.content", "p.comments_allowed", "p.created_at", "c.id", "c.post_id", "c.user_id", "c.parent_comment_id", "c.content", "c.created_at").
		From("posts p").
		LeftJoin("comments c on p.id = c.post_id").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db.DB).
		Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post

	var comment model.Comment

	commentMap := make(map[string]*model.Comment)

	postMap := make(map[string]*model.Post)

	for rows.Next() {
		var post model.Post

		if err = rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.AllowComments, &post.CreatedAt, &comment.ID, &comment.PostID, &comment.Author, &comment.ParentID, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}

		if _, exists := postMap[post.ID]; !exists {
			postMap[post.ID] = &post
			posts = append(posts, &post)
		}

		commentMap[comment.ID] = &comment
	}

	for _, post := range posts {
		post.Comments = buildCommentsTree(commentMap)
	}

	if rows.Err() != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, ErrNotFound
	}

	return posts, nil
}

func (r *Repository) CreatePost(ctx context.Context, post *model.Post) (string, error) {
	var id string

	userId, err := strconv.Atoi(post.Author)
	if err != nil {
		return "", err
	}

	err = sq.Insert("posts").
		Columns("user_id", "title", "content", "comments_allowed", "created_at").
		Values(userId, post.Title, post.Content, post.AllowComments, post.CreatedAt).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db.DB).
		QueryRow().
		Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *Repository) GetPostById(ctx context.Context, id string) (*model.Post, error) {
	var post model.Post

	err := sq.Select("id", "user_id", "title", "content", "comments_allowed", "created_at").
		From("posts").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db.DB).
		QueryRow().
		Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.AllowComments, &post.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWrongPostId
		}
		return nil, err
	}

	return &post, nil
}

func (r *Repository) GetPostByIdWithComments(ctx context.Context, id string) (*model.Post, error) {
	rows, err := sq.Select("p.id", "p.user_id", "p.title", "p.content", "p.comments_allowed", "p.created_at", "c.id", "c.post_id", "c.user_id", "c.parent_comment_id", "c.content", "c.created_at").
		From("posts p").
		LeftJoin("comments c on p.id = c.post_id").
		Where(sq.Eq{"p.id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db.DB).
		Query()

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrWrongPostId
		}

		return nil, err
	}
	defer rows.Close()

	var post model.Post
	var comment model.Comment

	commentMap := make(map[string]*model.Comment)

	for rows.Next() {

		if err = rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.AllowComments, &post.CreatedAt, &comment.ID, &comment.PostID, &comment.Author, &comment.ParentID, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}

		commentMap[comment.ID] = &comment
	}

	post.Comments = buildCommentsTree(commentMap)

	if rows.Err() != nil {
		return nil, err
	}

	return &post, nil
}

func (r *Repository) CreateComment(ctx context.Context, comment *model.Comment) (string, error) {
	var id string

	postId, err := strconv.Atoi(comment.PostID)
	if err != nil {
		return "", err
	}

	userId, err := strconv.Atoi(comment.Author)
	if err != nil {
		return "", err
	}

	values := []interface{}{postId, userId, comment.Content, comment.CreatedAt}
	columns := []string{"post_id", "user_id", "content", "created_at"}

	if comment.ParentID != nil {
		columns = append(columns, "parent_comment_id")
		values = append(values, *comment.ParentID)
	}

	err = sq.Insert("comments").
		Columns(columns...).
		Values(values...).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db.DB).
		QueryRow().
		Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}

func buildCommentsTree(comments map[string]*model.Comment) []*model.Comment {
	parentComments := make([]*model.Comment, 0)

	for _, comment := range comments {
		if *comment.ParentID == "" {
			parentComments = append(parentComments, comment)
		} else {
			parent := comments[*comment.ParentID]
			parent.Replies = append(parent.Replies, comment)
		}
	}

	return parentComments
}
