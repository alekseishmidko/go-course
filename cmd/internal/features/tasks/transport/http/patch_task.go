package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
	core_logger "github.com/alekseishmidko/go-course/cmd/internal/core/logger"
	core_http_request "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/request"
	core_http_response "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/response"
	core_http_types "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/types"
)

type PatchRequestDto struct {
	Title       core_http_types.Nullable[string] `json:"title" example:"Изучить Go" minLength:"1" maxLength:"100" swaggertype:"string"`
	Description core_http_types.Nullable[string] `json:"description" example:"Прочитать раздел о HTTP" minLength:"1" maxLength:"1000" swaggertype:"string"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" example:"true" swaggertype:"boolean"`
}

type PatchTaskResponse TaskDtoResponse

func (r *PatchRequestDto) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("title cannot be NULL")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` must be between 1 and 100 symbols")
		}
	}
	if r.Description.Set {
		if r.Description.Value == nil {
			return fmt.Errorf("description cannot be NULL")
		}

		descriptionLen := len([]rune(*r.Description.Value))
		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf("`Description` must be between 1 and 1000 symbols")
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("Completed cannot be null")
		}
	}

	return nil
}

// PatchTask godoc
// @Summary Обновить задачу
// @Tags tasks
// @Produce json
// @Param taskId path int true "ID задачи"
// @Accept json
// @Param request body PatchRequestDto true "Поля для обновления; непереданные поля сохраняются"
// @Description Пустой объект допустим. null запрещён для title, description и completed. completed обновляет completed_at автоматически.
// @Success 200 {object} PatchTaskResponse
// @Failure 400 {object} core_http_response.ErrorResponse
// @Failure 404 {object} core_http_response.ErrorResponse
// @Failure 409 {object} core_http_response.ErrorResponse
// @Failure 500 {object} core_http_response.ErrorResponse
// @Router /tasks/{taskId} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	taskId, err := core_http_request.GetIntPathValue(r, "taskId")

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get pathValue: taskId")
		return
	}
	var request PatchRequestDto
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode patch_task dto request")
		return
	}

	taskPatch := taskPatchFromRequest(request)
	taskDomain, err := h.tasksService.PatchTask(ctx, taskId, taskPatch)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	response := PatchTaskResponse(taskDtoFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}
func taskPatchFromRequest(dto PatchRequestDto) domain.TaskPatch {
	return domain.NewTaskPatch(dto.Title.ToDomain(), dto.Description.ToDomain(), dto.Completed.ToDomain())
}
