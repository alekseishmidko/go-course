package domain

import (
	"fmt"
	"time"

	core_errors "github.com/alekseishmidko/go-course/cmd/internal/core/error"
)

type Task struct {
	ID      int
	Version int

	Title       string
	Description *string
	Completed   bool

	CreatedAt   time.Time
	CompletedAt *time.Time

	AuthorUserId int
}

func NewTask(id int,
	version int,
	title string,
	description *string,
	completed bool,
	createdAt time.Time,
	completedAt *time.Time,
	authorUserId int) Task {

	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		AuthorUserId: authorUserId,
	}
}

func NewTaskUnitilazed(
	title string,
	description *string,
	authorUserId int,
) Task {
	return NewTask(UninitializedId, UninitializedVersion, title, description, false, time.Now(), nil, authorUserId)
}

func (t *Task) Validate() error {
	titleLength := len([]rune(t.Title))

	if titleLength < 1 || titleLength > 100 {
		return fmt.Errorf("invalid `Title` length: %d: %w", titleLength, core_errors.ErrInvalidArgument)
	}

	if t.Description != nil {
		descriptionLength := len([]rune(*t.Description))
		if descriptionLength < 1 || descriptionLength > 1000 {
			return fmt.Errorf("invalid `Description` length: %d: %w", descriptionLength, core_errors.ErrInvalidArgument)
		}
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf("`CompletedAt` cannot be nil if `Completed == true`: %w", core_errors.ErrInvalidArgument)
		}

		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf("`CompletedAt` can not be before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf("`CompletedAt` can not be nil if `Completed == false`: %w", core_errors.ErrInvalidArgument)
		}
	}

	return nil
}
