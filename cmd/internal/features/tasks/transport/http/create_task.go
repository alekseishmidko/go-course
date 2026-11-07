package tasks_transport_http

import (
	"net/http"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
	core_logger "github.com/alekseishmidko/go-course/cmd/internal/core/logger"
	core_http_request "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/request"
	core_http_response "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" example:"Изучить Go" minLength:"1" maxLength:"100" validate:"required,min=1,max=100"`
	Description  *string `json:"description" example:"Прочитать раздел о HTTP" minLength:"1" maxLength:"1000" extensions:"x-nullable" validate:"omitempty,min=1,max=1000"`
	AuthorUserId int     `json:"author_user_id" example:"1" validate:"required"`
}

type CreateTaskResponse TaskDtoResponse

// CreateTask godoc
// @Summary Создать задачу
// @Tags tasks
// @Produce json
// @Accept json
// @Param request body CreateTaskRequest true "Поля новой записи"
// @Success 201 {object} CreateTaskResponse
// @Failure 400 {object} core_http_response.ErrorResponse
// @Failure 500 {object} core_http_response.ErrorResponse
// @Router /tasks [post]
func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}
	taskDomain := domain.NewTaskUnitilazed(request.Title, request.Description, request.AuthorUserId)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create task")
		return
	}

	response := CreateTaskResponse(taskDtoFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func taskDtoFromDomain(task domain.Task) TaskDtoResponse {
	return TaskDtoResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserId: task.AuthorUserId,
	}
}
