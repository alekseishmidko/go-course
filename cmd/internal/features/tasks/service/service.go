package tasks_service

import (
	"context"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
)

type TasksService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTasks(ctx context.Context, userId *int, limit *int, offset *int) ([]domain.Task, error)
}

func NewTaskService(tasksRepository TasksRepository) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
	}
}
