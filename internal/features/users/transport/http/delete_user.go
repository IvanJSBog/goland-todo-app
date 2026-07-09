package users_transport_http

import (
	"net/http"

	core_logger "github.com/IvanJSBog/goland-todo-app/internal/core/logger"
	core_http_request "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/response"
)

// @DeleteUser godoc
// @Summary Удаление пользователя
// @Description Удаление существующего в системе пользователя по id
// @Tags users
// @Param id path int true "ID удаляемого пользователя"
// @Success 204 "успешное удаление пользователя"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 404 {object} core_http_response.ErrorResponse "User Not Found"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users/{id} [delete]
func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user id")
		return
	}

	err = h.usersService.DeleteUser(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")
	}

	responseHandler.NoContentResponse()
}
