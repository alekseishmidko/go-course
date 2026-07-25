package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/alekseishmidko/go-course/cmd/internal/core/error"
)

func (r *UsersRepository) DeleteUser(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.users
	WHERE id=$1`

	cmndTag, err := r.pool.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmndTag.RowsAffected() == 0 {

		return fmt.Errorf("user with id: %d:%w not found", id, core_errors.ErrNotFound)
	}

	return nil
}
