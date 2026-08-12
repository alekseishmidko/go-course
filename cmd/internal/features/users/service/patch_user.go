package users_service

import (
	"context"
	"fmt"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
)

func (s *UsersService) PatchUser(ctx context.Context, id int, patch domain.UserPatch) (domain.User, error) {
	// get by id
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user :%w", err)
	}
	// apply patch

	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf("error patching user: %w", err)
	}
	patchedUser, err := s.usersRepository.PatchUser(ctx, id, user)
	if err != nil {
		fmt.Errorf("patch user: %w", err)
	}
	return patchedUser, err
}
