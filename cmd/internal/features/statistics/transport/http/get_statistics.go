package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
	core_logger "github.com/alekseishmidko/go-course/cmd/internal/core/logger"
	core_http_request "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/request"
	core_http_response "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated            int      `json:"tasks_created"`
	TasksCompleted          int      `json:"tasks_completed"`
	TaskCompletedPercent    *float64 `json:"tasks_completed_percent"`
	TaskCompleteAverageTime *string  `json:"task_complete_average_time"`
}

func (h *StatisticHttpHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userId, from, to, err := getUserIdFromToQueryParams(r)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId/from/to query params")
		return
	}

	statistics, err := h.statisticService.GetStatistics(ctx, userId, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get statistics")
		return
	}
	response := toDtoFromDomain(statistics)
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIdFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIdParamKey = "userId"
		fromParamKey   = "from"
		toParamKey     = "to"
	)
	userId, err := core_http_request.GetIntQueryParam(r, userIdParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get userId query key:%w", err)
	}

	from, err := core_http_request.GetDateQueryParam(r, fromParamKey)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get from query key:%w", err)
	}

	to, err := core_http_request.GetDateQueryParam(r, toParamKey)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get to query key:%w", err)
	}

	return userId, from, to, nil
}
func toDtoFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TaskCompleteAverageTime != nil {
		duration := statistics.TaskCompleteAverageTime.String()
		avgTime = &duration
	}
	return GetStatisticsResponse{
		TasksCreated:            statistics.TasksCreated,
		TasksCompleted:          statistics.TasksCompleted,
		TaskCompletedPercent:    statistics.TaskCompletedPercent,
		TaskCompleteAverageTime: avgTime,
	}
}
