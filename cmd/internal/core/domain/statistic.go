package domain

import "time"

type Statistics struct {
	TasksCreated            int
	TasksCompleted          int
	TaskCompletedPercent    *float64
	TaskCompleteAverageTime *time.Duration
}

func NewStatistics(taskCreated int, tasksCompleted int, taskCompletedPercent *float64, taskCompletedAverageTime *time.Duration) Statistics {
	return Statistics{
		TasksCreated:            taskCreated,
		TasksCompleted:          tasksCompleted,
		TaskCompletedPercent:    taskCompletedPercent,
		TaskCompleteAverageTime: taskCompletedAverageTime,
	}
}
