package users_transport_http

import "github.com/IvanJSBog/goland-todo-app/internal/core/domain"

type UserDtoResponse struct {
	ID          int     `json:"id"`
	Version     int     `json:"version"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func userDTOFromDomain(user domain.User) UserDtoResponse {
	return UserDtoResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func usersDTOFromDomain(users []domain.User) []UserDtoResponse {
	usersDTO := make([]UserDtoResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO
}
