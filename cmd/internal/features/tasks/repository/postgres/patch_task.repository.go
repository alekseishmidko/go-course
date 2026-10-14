package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
	core_errors "github.com/alekseishmidko/go-course/cmd/internal/core/error"
	core_postgres_pool "github.com/alekseishmidko/go-course/cmd/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) PatchTask(ctx context.Context, taskId int, task domain.Task) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.tasks
	SET
	title=$1,
	description=$2,
	completed=$3,
	completed_at=$4,
	version=version + 1
	
	WHERE id=$5 AND version=$6
	
	RETURNING 
	id,version, title, description, created_at, completed, completed_at, author_user_id;`

	row := r.pool.QueryRow(ctx, query, task.Title, task.Description, task.Completed, task.CompletedAt, taskId, task.Version)
	var taskModel TaskModel

	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.CreatedAt,
		&taskModel.Completed,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserId,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf("task with id='%d' concurently accessed: %w", taskId, core_errors.ErrConflict)

		}
		return domain.Task{}, fmt.Errorf("scan task error: %w", err)
	}

	taskDomain := TaskDomainFromModel(taskModel)

	return taskDomain, nil
}
