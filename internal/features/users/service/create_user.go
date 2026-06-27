package users_service

import (
	"context"
	"fmt"

	"github.com/IvanJSBog/goland-todo-app/internal/core/domain"
)

func (u *UsersService) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	err := user.Validate()
	if err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	user, err = u.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
