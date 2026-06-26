package users_postgres_repository

import "github.com/IvanJSBog/goland-todo-app/internal/core/domain"

type UsersModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func userDomainsFromModels(models []UsersModel) []domain.User {
	userDomains := make([]domain.User, len(models))

	for i, model := range models {
		userDomains[i] = domain.NewUser(model.ID, model.Version, model.FullName, model.PhoneNumber)
	}

	return userDomains
}
