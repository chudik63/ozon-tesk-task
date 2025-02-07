package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Config struct {
	UserName string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	DbName   string `env:"POSTGRES_DB"`
}

type DB struct {
	*sql.DB
}

func New(config Config) (DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s", config.UserName, config.Password, config.DbName, config.Host, config.Port)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return DB{}, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return DB{}, fmt.Errorf("failed to ping database: %w", err)
	}

	return DB{db}, nil
}
