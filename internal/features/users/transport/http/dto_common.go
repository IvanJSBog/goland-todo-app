package users_transport_http

import "github.com/IvanJSBog/goland-todo-app/internal/core/domain"

type UserDtoResponse struct {
	ID          int     `json:"id" example:"1"`
	Version     int     `json:"version" example:"1"`
	FullName    string  `json:"full_name" example:"John Doe"`
	PhoneNumber *string `json:"phone_number" example:"+792122822890"`
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
