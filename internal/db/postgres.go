package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"crud_service/internal/config"
)

func NewPool(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	//pool, err := pgxpool.New(ctx, cfg.DSN())

	config, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.MaxConns = 100                      // максимум соединений
	config.MinConns = 5                        // держать минимум открытыми
	config.MaxConnLifetime = 1 * time.Hour     // максимальное время жизни соединения
	config.MaxConnIdleTime = 30 * time.Minute  // закрыть если простаивает 30 минут
	config.HealthCheckPeriod = 1 * time.Minute // проверять живые соединения

	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		return nil, fmt.Errorf("Unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return pool, fmt.Errorf("Unable to ping database: %w", err)
	}

	return pool, nil
}
