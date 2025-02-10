package repository

import (
	"context"
	"database/sql"
	"errors"
	"ozon-tesk-task/internal/database"
	"ozon-tesk-task/internal/transport/graph/model"
	"ozon-tesk-task/pkg/pointer"

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
		OrderBy("p.created_at, c.created_at").
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

	commentMap := make(map[int32][]*model.Comment)

	postMap := make(map[int32]*model.Post)

	for rows.Next() {
		var (
			post      model.Post
			comment   model.Comment
			id        sql.NullInt32
			postId    sql.NullInt32
			author    sql.NullInt32
			parentId  sql.NullInt32
			content   sql.NullString
			createdAt sql.NullString
		)

		if err = rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.AllowComments, &post.CreatedAt, &id, &postId, &author, &parentId, &content, &createdAt); err != nil {
			return nil, err
		}

		if id.Valid {
			comment = model.Comment{
				ID:        id.Int32,
				PostID:    postId.Int32,
				Author:    author.Int32,
				ParentID:  &parentId.Int32,
				Content:   content.String,
				CreatedAt: createdAt.String,
			}

			commentMap[post.ID] = append(commentMap[post.ID], &comment)
		}

		if _, exists := postMap[post.ID]; !exists {
			postMap[post.ID] = &post
			posts = append(posts, &post)
		}
	}

	for _, post := range posts {
		post.Comments = buildCommentsTree(commentMap[post.ID])
	}

	if rows.Err() != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, ErrNotFound
	}

	return posts, nil
}

func (r *Repository) CreatePost(ctx context.Context, post *model.Post) (int32, error) {
	var id int32

	err := sq.Insert("posts").
		Columns("user_id", "title", "content", "comments_allowed", "created_at").
		Values(post.Author, post.Title, post.Content, post.AllowComments, post.CreatedAt).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db.DB).
		QueryRow().
		Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) GetPostById(ctx context.Context, id int32) (*model.Post, error) {
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

func (r *Repository) GetPostByIdWithComments(ctx context.Context, id int32) (*model.Post, error) {
	rows, err := sq.Select("p.id", "p.user_id", "p.title", "p.content", "p.comments_allowed", "p.created_at", "c.id", "c.post_id", "c.user_id", "c.parent_comment_id", "c.content", "c.created_at").
		From("posts p").
		LeftJoin("comments c on p.id = c.post_id").
		Where(sq.Eq{"p.id": id}).
		OrderBy("c.created_at").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db.DB).
		Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var post model.Post

	comments := make([]*model.Comment, 0)

	if !rows.Next() {
		return nil, ErrWrongPostId
	}

	for {
		var (
			comment   model.Comment
			id        sql.NullInt32
			postId    sql.NullInt32
			author    sql.NullInt32
			parentId  sql.NullInt32
			content   sql.NullString
			createdAt sql.NullString
		)

		if err = rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.AllowComments, &post.CreatedAt, &id, &postId, &author, &parentId, &content, &createdAt); err != nil {
			return nil, err
		}

		if id.Valid {
			comment = model.Comment{
				ID:        id.Int32,
				PostID:    postId.Int32,
				Author:    author.Int32,
				ParentID:  &parentId.Int32,
				Content:   content.String,
				CreatedAt: createdAt.String,
			}

			comments = append(comments, &comment)
		}

		if !rows.Next() {
			break
		}
	}

	post.Comments = buildCommentsTree(comments)

	if rows.Err() != nil {
		return nil, err
	}

	return &post, nil
}

func (r *Repository) CreateComment(ctx context.Context, comment *model.Comment) (int32, error) {
	var id int32

	values := []interface{}{comment.PostID, comment.Author, comment.Content, comment.CreatedAt}
	columns := []string{"post_id", "user_id", "content", "created_at"}

	if comment.ParentID != nil {
		columns = append(columns, "parent_comment_id")
		values = append(values, *comment.ParentID)
	}

	err := sq.Insert("comments").
		Columns(columns...).
		Values(values...).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db.DB).
		QueryRow().
		Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) GetCommentById(ctx context.Context, commentId int32) (*model.Comment, error) {
	var comment model.Comment

	err := sq.Select("id", "post_id", "user_id", "parent_comment_id", "content", "created_at").
		From("comments").
		Where(sq.Eq{"id": commentId}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db.DB).
		QueryRow().
		Scan(&comment.ID, &comment.PostID, &comment.Author, &comment.ParentID, &comment.Content, &comment.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWrongCommentId
		}

		return nil, err
	}

	return &comment, nil
}

func buildCommentsTree(comments []*model.Comment) []*model.Comment {
	commentMap := make(map[int32]*model.Comment)
	parentComments := make([]*model.Comment, 0)

	for _, comment := range comments {
		commentMap[comment.ID] = comment
	}

	for _, comment := range comments {
		if pointer.Deref(comment.ParentID, 0) == 0 {
			parentComments = append(parentComments, comment)
		} else {
			parent := commentMap[*comment.ParentID]
			parent.Replies = append(parent.Replies, comment)
		}
	}

	return parentComments
}
