package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/alekseishmidko/go-course/cmd/internal/core/logger"
	core_http_request "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/request"
	core_http_response "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/response"
)

type GetTaskResponse TaskDtoResponse

func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(logger, rw)
	taskId, err := core_http_request.GetIntPathValue(r, "taskId")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get pathValue: taskId")
		return
	}
	taskDomain, err := h.tasksService.GetTask(ctx, taskId)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task")
		return
	}
	response := GetTaskResponse(taskDtoFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}
