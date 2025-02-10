package repository

import "errors"

var (
	ErrNotFound             = errors.New("nothing was found")
	ErrWrongPostId          = errors.New("post with such id does not exist")
	ErrWrongCommentId       = errors.New("comment with such id does not exist")
	ErrCommentsNotAllowed   = errors.New("post with such id does not allow comments")
	ErrMatchCommentWithPost = errors.New("comment with such id does not belong to the post")
)
