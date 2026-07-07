package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	// pool, err := pgxpool.New(context.Background(), dsn)
	// if err != nil {
	// 	return nil, fmt.Errorf("connection pool failed:%w", err)
	// }
	// defer pool.Close()
	// if err := pool.Ping(context.Background()); err != nil {
	// 	return nil, fmt.Errorf("connection pool ping failed:%w", err)
	// }
	// return pool, nil

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("creating pgxpool config failed: %w", err)
	}

	config.MinConns = 5
	config.MaxConns = 25
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 15 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute
	config.ConnConfig.ConnectTimeout = 5 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("creating pgxpool failed: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("connection pool ping failed: %w", err)
	}
	return pool, nil
}
