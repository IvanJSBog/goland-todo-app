package users_service

import (
	"context"
	"fmt"

	"github.com/IvanJSBog/goland-todo-app/internal/core/domain"
)

func (u *UsersService) PatchUser(
	ctx context.Context,
	id int,
	patch domain.PatchUser,
) (domain.User, error) {
	user, err := u.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user: %w", err)
	}

	err = user.ApplyPatch(patch)
	if err != nil {
		return domain.User{}, fmt.Errorf("apply user patch: %w", err)
	}
	patchedUser, err := u.usersRepository.PatchUser(ctx, id, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("patch user: %w", err)
	}
	return patchedUser, nil
}
