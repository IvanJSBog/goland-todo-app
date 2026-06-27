package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/IvanJSBog/goland-todo-app/internal/core/domain"
)

func (u *UsersRepository) GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, full_name, phone_number 
		FROM todoapp.users
		ORDER BY id ASC
		LIMIT $1
		OFFSET $2;
	`

	rows, err := u.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select users %w", err)
	}
	defer rows.Close()

	var userModels []UsersModel
	for rows.Next() {
		var userModel UsersModel

		err := rows.Scan(&userModel.ID, &userModel.Version, &userModel.FullName, &userModel.PhoneNumber)
		if err != nil {
			return nil, fmt.Errorf("scan user model %w", err)
		}

		userModels = append(userModels, userModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows %w", err)
	}

	userDomains := userDomainsFromModels(userModels)
	
	return userDomains, nil
}
