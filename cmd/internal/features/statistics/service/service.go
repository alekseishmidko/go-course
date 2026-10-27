package statistics_service

import (
	"context"
	"time"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetTasks(ctx context.Context, userId *int, from *time.Time, to *time.Time) ([]domain.Task, error)
}

func NewStatisticsService(statisticsRepository StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: statisticsRepository,
	}
}
