package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func ConnectPostgres(dbURL string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	return conn, nil
}
