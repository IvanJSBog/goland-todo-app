package users_transport_http

import (
	"net/http"

	core_logger "github.com/IvanJSBog/goland-todo-app/internal/core/logger"
	core_http_response "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/response"
	core_http_utils "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/utils"
)

func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	id, err := core_http_utils.GetIntPathValue(r, "id")
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
