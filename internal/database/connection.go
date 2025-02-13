package database

import (
	"context"
	"errors"
	"fmt"
	"ozon-tesk-task/internal/config"
	"ozon-tesk-task/pkg/logger"

	"go.uber.org/zap"
)

func NewDatabase(ctx context.Context, cfg *config.Config) (*Database, error) {
	var (
		db    *Database
		err   error
		dsn   string
		dbURL string
	)

	l := logger.GetLoggerFromCtx(ctx)

	switch cfg.StorageType {
	case "postgres":
		db = New(cfg, "postgres")
		dsn = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s", cfg.UserName, cfg.Password, cfg.DbName, cfg.Host, cfg.Port)
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)

		l.Debug(ctx, "Connecting to postgres database", zap.String("dsn", dsn), zap.String("dbURL", dbURL))
	case "memory":
		db = New(cfg, "sqlite")
		dsn = "file::memory:?cache=shared"
		dbURL = "sqlite://file::memory:?cache=shared"

		l.Debug(ctx, "Connecting to sqlite database", zap.String("dsn", dsn), zap.String("dbURL", dbURL))
	default:
		return nil, errors.New("invalid storage type")
	}

	err = db.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	err = db.MigrateUp(ctx, dbURL)
	if err == MigrationNoChange {
		l.Info(ctx, "No change after migration")
	} else if err != nil {
		return nil, err
	}

	l.Info(ctx, "Migration completed succesfully")

	return db, nil
}
