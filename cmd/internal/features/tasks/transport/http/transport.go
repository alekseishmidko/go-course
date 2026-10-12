package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
	core_http_server "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/server"
)

type TasksHTTPHandler struct {
	tasksService TasksService
}

type TasksService interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTasks(ctx context.Context, userId *int, limit *int, offset *int) ([]domain.Task, error)
	GetTask(ctx context.Context, taskId *int) (domain.Task, error)
	DeleteTask(ctx context.Context, taskId int) error
}

func NewTaskHttpHandler(taskService TasksService) *TasksHTTPHandler {
	return &TasksHTTPHandler{tasksService: taskService}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{Method: http.MethodPost, Path: "/tasks", Handler: h.CreateTask},
		{Method: http.MethodGet, Path: "/tasks", Handler: h.GetTasks},
		{Method: http.MethodGet, Path: "/tasks/{taskId}", Handler: h.GetTask},
		{Method: http.MethodDelete, Path: "/tasks/{taskId}", Handler: h.DeleteTask},
	}
}
