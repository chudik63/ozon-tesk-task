package repository

import "ozon-tesk-task/internal/database/sql"

type Repository struct {
	db *sql.Database
}

func New(db *sql.Database) *Repository {
	return &Repository{db: db}
}
