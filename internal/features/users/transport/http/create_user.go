package users_transport_http

import (
	"net/http"

	"github.com/IvanJSBog/goland-todo-app/internal/core/domain"
	core_logger "github.com/IvanJSBog/goland-todo-app/internal/core/logger"
	core_http_request "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100" example:"John Doe"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+" example:"+792122822890"`
}

type CreateUserResponse UserDtoResponse

// CreateUser godoc
// @Summary Создать пользователя
// @Description Создать нового пользователя в системе
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "CreateUser тело запроса"
// @Success 201 {object} CreateUserResponse "Успешно созданный пользователь"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users [post]
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
