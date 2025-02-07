package postgres

import (
	"context"
	"fmt"
	"ozon-tesk-task/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(ctx context.Context, cfg Config, source string) error {
	l := logger.GetLoggerFromCtx(ctx)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)

	m, err := migrate.New("file://"+source, dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		l.Info(ctx, "no change after migration")
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to make migration up: %w", err)
	}

	l.Info(ctx, "migration complete succesfully")

	return nil
}
