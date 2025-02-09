package database

import (
	"context"
	"database/sql"
	"fmt"
	"ozon-tesk-task/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

var (
	MigrationNoChange = migrate.ErrNoChange
)

type Database struct {
	DB     *sql.DB
	driver string
	config *config.Config
}

func New(cfg *config.Config, driver string) *Database {
	return &Database{
		config: cfg,
		driver: driver,
	}
}

func (d *Database) Connect(ctx context.Context, dsn string) error {
	db, err := sql.Open(d.driver, dsn)

	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	d.DB = db

	return nil
}

func (d *Database) MigrateUp(ctx context.Context, dbURL string) error {
	m, err := migrate.New(fmt.Sprintf("file://%s/%s", d.config.MigrationsPath, d.driver), dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		return MigrationNoChange
	}

	if err != nil {
		return fmt.Errorf("failed to make migration up: %w", err)
	}

	return nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}
