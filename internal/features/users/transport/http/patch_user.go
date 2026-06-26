package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/IvanJSBog/goland-todo-app/internal/core/domain"
	core_logger "github.com/IvanJSBog/goland-todo-app/internal/core/logger"
	core_http_request "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/response"
	core_http_types "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/types"
	core_http_utils "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/utils"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (h *PatchUserRequest) Validate() error {
	if h.FullName.Set {
		if h.FullName.Value == nil {
			return fmt.Errorf("FullName cannot be null")
		}

		fullNameLen := len([]rune(*h.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("FullName length must be between 3 and 100")
		}
	}

	if h.PhoneNumber.Set {
		if h.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*h.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("PhoneNumber length must be between 10 and 15")
			}

			if !strings.HasPrefix(*h.PhoneNumber.Value, "+") {
				return fmt.Errorf("PhoneNumber must start with '+'")
			}
		}
	}

	return nil
}

type PatchUserResponse UserDtoResponse

func (u *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userId, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)
		return
	}

	var body PatchUserRequest
	err = core_http_request.DecodeAndValidateRequest(r, &body)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")
		return
	}
	userPatch := UserPatchFromRequest(body)
	userDomain, err := u.usersService.PatchUser(ctx, userId, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func UserPatchFromRequest(request PatchUserRequest) domain.PatchUser {
	return domain.PatchUser{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}
