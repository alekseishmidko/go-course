package core_pgx_pool

import (
	"errors"
	"testing"

	core_postgres_pool "github.com/alekseishmidko/go-course/cmd/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type stubRow struct {
	err       error
	scanCalls int
	dest      []any
}

func (r *stubRow) Scan(dest ...any) error {
	r.scanCalls++
	r.dest = dest
	return r.err
}

func TestPgxRowScanDelegatesToWrappedRow(t *testing.T) {
	wrapped := &stubRow{}
	row := pgxRow{Row: wrapped}
	var id int

	if err := row.Scan(&id); err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if wrapped.scanCalls != 1 {
		t.Fatalf("wrapped Scan() calls = %d, want 1", wrapped.scanCalls)
	}
	if len(wrapped.dest) != 1 || wrapped.dest[0] != &id {
		t.Fatal("Scan() did not forward destinations to the wrapped row")
	}
}

func TestPgxRowScanMapsNoRows(t *testing.T) {
	row := pgxRow{Row: &stubRow{err: pgx.ErrNoRows}}

	if err := row.Scan(); !errors.Is(err, core_postgres_pool.ErrNoRows) {
		t.Fatalf("Scan() error = %v, want ErrNoRows", err)
	}
}

func TestPgxRowScanMapsForeignKeyViolation(t *testing.T) {
	row := pgxRow{Row: &stubRow{err: &pgconn.PgError{Code: "23503"}}}

	if err := row.Scan(); !errors.Is(err, core_postgres_pool.ErrWiolatesForeignKey) {
		t.Fatalf("Scan() error = %v, want ErrWiolatesForeignKey", err)
	}
}

func TestPgxRowScanPreservesUnknownError(t *testing.T) {
	original := errors.New("scan failed")
	row := pgxRow{Row: &stubRow{err: original}}

	if err := row.Scan(); !errors.Is(err, original) {
		t.Fatalf("Scan() error = %v, want original error", err)
	}
}
