package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/IvanJSBog/goland-todo-app/internal/core/logger"
	core_http_request "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/response"
)

type GetUsersResponse []UserDtoResponse

// GetUsers godoc
// @Summary Список пользователей
// @Description Просмотр списка пользователей с опциональной пагинацией
// @Tags users
// @Produce json
// @Param limit query int false "Размер страницы с пользователями"
// @Param offset query int false "Смещение страницы с пользователями"
// @Success 200 {object} GetUsersResponse "Успешное получение списка пользователей"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users [get]
func (u *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := getLimitOffsetQueryParams(req)
	if err != nil {
		responseHandler.ErrorResponse(
			err, "failed to get limit and offset query params",
		)

		return
	}

	userDomains, err := u.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err, "failed to get users",
		)

		return
	}

	response := GetUsersResponse(usersDTOFromDomain(userDomains))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := core_http_request.GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get limit parameter: %w", err)
	}
	offset, err := core_http_request.GetIntQueryParam(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get offset parameter: %w", err)
	}

	return limit, offset, nil
}
