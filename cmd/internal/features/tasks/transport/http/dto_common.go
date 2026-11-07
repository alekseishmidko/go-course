package tasks_transport_http

import (
	"time"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
)

type TaskDtoResponse struct {
	ID      int `json:"id"`
	Version int `json:"version"`

	Title       string  `json:"title"`
	Description *string `json:"description" extensions:"x-nullable"`
	Completed   bool    `json:"completed"`

	CreatedAt   time.Time  `json:"created_at" format:"date-time"`
	CompletedAt *time.Time `json:"completed_at" format:"date-time" extensions:"x-nullable"`

	AuthorUserId int `json:"author_user_id"`
}

func taskDtosFromDomains(tasks []domain.Task) []TaskDtoResponse {
	dtos := make([]TaskDtoResponse, len(tasks))
	for i, task := range tasks {
		dtos[i] = taskDtoFromDomain(task)
	}

	return dtos
}
