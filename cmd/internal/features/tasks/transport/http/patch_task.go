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
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"`
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
