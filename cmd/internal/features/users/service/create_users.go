package users_service

import (
	"context"
	"fmt"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
)

func (service *UsersService) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("user validation failed: %w", err)
	}

	user, err := service.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
