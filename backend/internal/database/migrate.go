package database

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Up(ctx context.Context, pool *pgxpool.Pool) error {
	return run(ctx, pool, func(m *migrate.Migrate) error {
		err := m.Up()
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	})
}

func Down(ctx context.Context, pool *pgxpool.Pool) error {
	return run(ctx, pool, func(m *migrate.Migrate) error {
		err := m.Steps(-1)
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	})
}

func run(ctx context.Context, pool *pgxpool.Pool, action func(*migrate.Migrate) error) error {
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("check postgres before migration: %w", err)
	}

	sqlDB := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create migration database driver: %w", err)
	}
	source, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("load embedded migrations: %w", err)
	}
	migration, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("create migration runner: %w", err)
	}
	defer migration.Close()

	if err := action(migration); err != nil {
		return fmt.Errorf("run migration: %w", err)
	}
	return nil
}
