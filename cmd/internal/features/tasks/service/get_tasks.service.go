package tasks_service

import (
	"context"
	"fmt"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
	core_errors "github.com/alekseishmidko/go-course/cmd/internal/core/error"
)

func (s *TasksService) GetTasks(ctx context.Context, userId *int, limit *int, offset *int) ([]domain.Task, error) {

	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("Limit can not be negative: %w", core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset can not be negative: %w", core_errors.ErrInvalidArgument)
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, userId, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("get tasks from repository: %w", err)
	}

	return tasks, nil
}
