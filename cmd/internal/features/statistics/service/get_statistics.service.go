package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
	core_errors "github.com/alekseishmidko/go-course/cmd/internal/core/error"
)

func (s *StatisticsService) GetStatistics(ctx context.Context, userId *int, from *time.Time, to *time.Time) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf("to must be after from: %w", core_errors.ErrInvalidArgument)
		}
	}

	tasks, err := s.statisticsRepository.GetTasks(ctx, userId, from, to)

	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get tasks from repository :%w", err)
	}

	statistics := calculateStatistics(tasks)

	return statistics, nil
}

func calculateStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.Statistics{TasksCreated: 0, TasksCompleted: 0, TaskCompletedPercent: nil, TaskCompleteAverageTime: nil}
	}

	tasksCreated := len(tasks)

	tasksCompleted := 0
	var taskCompleteAverageTime time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted = tasksCompleted + 1
		}
		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			taskCompleteAverageTime += *completionDuration
		}
	}

	taskCompletedPercent := float64(tasksCompleted) / float64(tasksCreated) * 100
	var taskAverageCompletionTime *time.Duration

	if tasksCompleted > 0 && taskCompleteAverageTime != 0 {
		avg := taskCompleteAverageTime / time.Duration(tasksCompleted)

		taskAverageCompletionTime = &avg
	}

	return domain.Statistics{TasksCreated: tasksCreated, TasksCompleted: tasksCompleted, TaskCompletedPercent: &taskCompletedPercent, TaskCompleteAverageTime: taskAverageCompletionTime}
}
