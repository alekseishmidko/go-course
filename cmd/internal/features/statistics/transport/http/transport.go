package statistics_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
	core_http_server "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/server"
)

type StatisticHttpHandler struct {
	statisticService StatisticService
}

type StatisticService interface {
	GetStatistics(ctx context.Context, userId *int, from *time.Time, to *time.Time) (domain.Statistics, error)
}

func NewStatisticHttpHandler(statisticService StatisticService) (h *StatisticHttpHandler) {

	return &StatisticHttpHandler{statisticService: statisticService}
}

func (h *StatisticHttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{{Method: http.MethodGet, Path: "/statistics", Handler: h.GetStatistics}}
}
