package db

import (
	"context"
	"crud_service/internal/config"
	"fmt"
	"net"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func NewConnect(cfg config.ClickHouseConfig) (driver.Conn, error) {
	ctx := context.Background()

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{net.JoinHostPort(cfg.Host, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.Dbname,
			Username: cfg.User,
			Password: cfg.Password,
		},
		MaxOpenConns: 5,
	})

	if err != nil {
		return nil, fmt.Errorf("clickhouse.NewConnect Open failed: %w", err)
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("clickhouse.NewConnect Ping failed: %w", err)
	}

	return conn, nil
}
