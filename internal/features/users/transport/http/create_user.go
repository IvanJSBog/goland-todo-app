package users_transport_http

import (
	"net/http"

	"github.com/IvanJSBog/goland-todo-app/internal/core/domain"
	core_logger "github.com/IvanJSBog/goland-todo-app/internal/core/logger"
	core_http_request "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"`
}

type CreateUserResponse UserDtoResponse

func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("Invoke CreateUser Handler")
	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(req, &request); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate")
		return
	}

	userDomain := domainFromDTO(request)

	createdUser, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to create user")
		return
	}

	response := CreateUserResponse(userDTOFromDomain(createdUser))

	responseHandler.JSONResponse(response, http.StatusCreated)

	rw.WriteHeader(http.StatusOK)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
