package storage

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"

	"github.com/exaring/otelpgx"
	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// New creates a new pgxpool.Pool, connects to PostgreSQL using the provided DSN,
// and runs embedded migrations before returning
func New(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("storage: parse dsn: %w", err)
	}
	poolCfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("storage: connect pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("storage: ping postgres: %w", err)
	}

	if err := runMigrations(pool, dsn); err != nil {
		pool.Close()
		return nil, fmt.Errorf("storage: run migrations: %w", err)
	}

	return pool, nil
}

// runMigrations embeds the SQL files and applies all pending migrations
func runMigrations(pool *pgxpool.Pool, dsn string) error {
	// Build iofs source from embedded filesystem.
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("migrations: sub fs: %w", err)
	}

	srcDriver, err := iofs.New(sub, ".")
	if err != nil {
		return fmt.Errorf("migrations: iofs source: %w", err)
	}
	defer srcDriver.Close()

	// Open a standard database/sql connection from the pgx pool for golang-migrate
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	dbDriver, err := migratepgx.WithInstance(db, &migratepgx.Config{})
	if err != nil {
		return fmt.Errorf("migrations: db driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", srcDriver, "pgx5", dbDriver)
	if err != nil {
		return fmt.Errorf("migrations: new instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			// Already up to date
			return nil
		}
		return fmt.Errorf("migrations: up: %w", err)
	}

	return nil
}
