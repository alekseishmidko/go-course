package tasks_service

import (
	"context"
	"fmt"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
)

func (s *TasksService) GetTask(ctx context.Context, taskId *int) (domain.Task, error) {

	task, err := s.tasksRepository.GetTask(ctx, taskId)

	if err != nil {
		return domain.Task{}, fmt.Errorf("get task from repository: %w", err)
	}

	return task, nil
}
