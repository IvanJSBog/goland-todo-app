package core_pgx_pool

import (
	"errors"

	core_postgres_pool "github.com/IvanJSBog/goland-todo-app/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgxRows struct {
	pgx.Rows
}

func (r pgxRows) Scan(dest ...any) error {
	err := r.Rows.Scan(dest...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_postgres_pool.ErrNoRows
		}
		return err
	}
	return nil
}

type pgxRow struct {
	pgx.Row
}

type pgxCommandTag struct {
	pgconn.CommandTag
}
