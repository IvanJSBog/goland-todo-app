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
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"John Doe"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+793135415151"`
}

// PatchUser godoc
// @Summary Изменение пользователя
// @Description Изменение информации об уже существующем в системе пользователе
// @Description ### Логика обновления полей (Three-state logic):
// @Description 1. **Поле не передано**: `phone_number` игнорируется, значение в БД не меняется
// @Description 2. **Явно передано значение**: `"phone_number":"+793134533451"` - устанавливается новый номер телефона в БД
// @Description 3. **Передан null**: `"phone_number" : null` - очищает поле в БД (set to NULL)
// @Description Ограничения `full_name` не может быть выставлен в null
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID изменяемого пользователя"
// @Param request body PatchUserRequest true "PatchUser тело запроса"
// @Success 200 {object} PatchUserResponse "Успешно измененный пользователь"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users/{id} [patch]
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

	userId, err := core_http_request.GetIntPathValue(r, "id")
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
	return domain.NewUserPatch(request.FullName.ToDomain(), request.PhoneNumber.ToDomain())
}
