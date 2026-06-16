package users_service

import (
	"context"
	"fmt"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
	core_errors "github.com/alekseishmidko/go-course/cmd/internal/core/error"
)

func (service *UsersService) GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("Limit can not be negative: %w", core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset can not be negative: %w", core_errors.ErrInvalidArgument)
	}
	users, err := service.usersRepository.GetUsers(ctx, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("get users from repository: %w", err)
	}

	return users, nil
}
