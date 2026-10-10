package tasks_transport_http

import (
	"time"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
)

type TaskDtoResponse struct {
	ID      int `json:"id"`
	Version int `json:"version"`

	Title       string  `json:"title"`
	Description *string `json:"description"`
	Completed   bool    `json:"completed"`

	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`

	AuthorUserId int `json:"author_user_id"`
}

func taskDtosFromDomains(tasks []domain.Task) []TaskDtoResponse {
	dtos := make([]TaskDtoResponse, len(tasks))
	for i, task := range tasks {
		dtos[i] = taskDtoFromDomain(task)
	}

	return dtos
}
