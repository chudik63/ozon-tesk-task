package repository

import "errors"

var (
	ErrNotFound    = errors.New("nothing was found")
	ErrWrongPostId = errors.New("post with such id does not exist")
)
