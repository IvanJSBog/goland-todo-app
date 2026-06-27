package core_pgx_pool

import (
	"context"
	"fmt"
	"time"

	core_postgres_pool "github.com/IvanJSBog/goland-todo-app/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxConnectionPool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func NewConnectionPool(ctx context.Context, config Config) (*PgxConnectionPool, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	pgxconfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse pgxconfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("new pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping pool: %w", err)
	}

	return &PgxConnectionPool{Pool: pool, opTimeout: config.Timeout}, nil
}

func (p *PgxConnectionPool) OpTimeout() time.Duration {
	return p.opTimeout
}

func (p *PgxConnectionPool) Query(ctx context.Context, sql string, args ...any) (core_postgres_pool.Rows, error) {
	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgxRows{rows}, nil
}
func (p *PgxConnectionPool) QueryRow(ctx context.Context, sql string, args ...any) core_postgres_pool.Row {
	row := p.Pool.QueryRow(ctx, sql, args...)
	return pgxRow{row}
}
func (p *PgxConnectionPool) Exec(ctx context.Context, sql string, arguments ...interface{}) (core_postgres_pool.CommandTag, error) {
	commandTag, err := p.Pool.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	return pgxCommandTag{
		commandTag,
	}, nil
}
