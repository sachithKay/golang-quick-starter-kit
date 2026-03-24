package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sachithKay/ghost/internal/config"
)

// NewPool is our "Database Factory" (Like a Spring Bean Provider).
// It returns a *pgxpool.Pool which handles concurrent connections for us.
func NewPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	// 1. Parse the connection string into a config object
	poolConfig, err := pgxpool.ParseConfig(cfg.DB.URL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	// 2. Apply our "Prod-Grade" settings from our Config struct
	poolConfig.MaxConns = int32(cfg.DB.MaxConnections)
	poolConfig.MinConns = int32(cfg.DB.MinConnections)
	poolConfig.MaxConnIdleTime = cfg.DB.ConnMaxIdleTime

	// 3. Create the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// 4. "Ping" the database to ensure the credentials/network are actually working.
	// In Go, creating a pool is "lazy"—it won't fail until you try to use it.
	// We use a timeout to ensure we don't hang forever if the DB is down.
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		return nil, fmt.Errorf("database heartbeat failed: %w", err)
	}

	return pool, nil
}