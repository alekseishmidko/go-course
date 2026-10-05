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
}

func NewTaskHttpHandler(taskService TasksService) *TasksHTTPHandler {
	return &TasksHTTPHandler{tasksService: taskService}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{Method: http.MethodPost, Path: "/tasks", Handler: h.CreateTask},
	}
}
