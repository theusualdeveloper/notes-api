package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	err := createSchemaTable(ctx, pool)
	if err != nil {
		return err
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading migration files failed: %w", err)
	}
	slices.SortFunc(files, func(a, b os.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})
	for _, file := range files {
		exists, err := existsInSchema(ctx, pool, file.Name())
		if err != nil {
			return err
		}
		if !exists {
			err = runTransaction(ctx, pool, dir, file.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createSchemaTable(ctx context.Context, pgxpool *pgxpool.Pool) error {
	q := `CREATE TABLE IF NOT EXISTS schema_migrations (
    id SERIAL PRIMARY KEY,
    version TEXT NOT NULL,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW())`
	_, err := pgxpool.Exec(ctx, q)
	if err != nil {
		return fmt.Errorf("creating schema table failed: %w", err)
	}
	return nil
}

func existsInSchema(ctx context.Context, pgxpool *pgxpool.Pool, fname string) (bool, error) {
	var exists bool
	q := `SELECT EXISTS (SELECT id FROM schema_migrations WHERE version = $1)`
	err := pgxpool.QueryRow(ctx, q, fname).Scan(&exists)
	if err != nil {
		return exists, fmt.Errorf("executing existence query failed: %w", err)
	}
	return exists, nil
}

func runTransaction(
	ctx context.Context,
	pgxpool *pgxpool.Pool,
	dir string,
	fname string,
) error {
	tx, err := pgxpool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction failed: %w", err)
	}
	defer tx.Rollback(ctx)
	err = runMigration(ctx, tx, dir, fname)
	if err != nil {
		return err
	}
	err = createMigration(ctx, tx, fname)
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing failed: %w", err)
	}
	return nil
}

func runMigration(
	ctx context.Context,
	tx pgx.Tx,
	dir string,
	fname string,
) error {
	fname = filepath.Join(dir, fname)
	qByte, err := os.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("reading file failed: %w", err)
	}
	_, err = tx.Exec(ctx, string(qByte))
	if err != nil {
		return fmt.Errorf("executing migration failed: %w", err)
	}
	return nil
}

func createMigration(
	ctx context.Context,
	tx pgx.Tx,
	fname string,
) error {
	q := `INSERT INTO schema_migrations (version) VALUES ($1)`
	_, err := tx.Exec(ctx, q, fname)
	if err != nil {
		return fmt.Errorf("insert into schema migrations failed: %w", err)
	}
	return nil
}
