package core_pgx_pool

import (
	"errors"
	"fmt"

	core_postgres_pool "github.com/alekseishmidko/go-course/cmd/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgxRows struct {
	pgx.Rows
}
type pgxRow struct {
	pgx.Row
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

func (r pgxRow) Scan(dest ...any) error {
	err := r.Scan()

	if err != nil {
		mapErrors(err)
	}
	return nil
}

func mapErrors(err error) error {
	const ViolatesForeignKeyCode = "23504"
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_postgres_pool.ErrNoRows
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {

			if pgErr.Code == ViolatesForeignKeyCode {

				return fmt.Errorf("%v: %w", err, core_postgres_pool.ErrWiolatesForeignKey)

			}

		}
		return err
	}
	return fmt.Errorf("%v: %w", err, core_postgres_pool.ErrUnknown)
}
