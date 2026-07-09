package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/IvanJSBog/goland-todo-app/internal/core/domain"
	core_logger "github.com/IvanJSBog/goland-todo-app/internal/core/logger"
	core_http_request "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/response"
	core_http_types "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title" swaggertype:"string" example:"Some Task"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"This is a description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" swaggertype:"boolean"`
}

type PatchTaskResponse TaskDTOResponse

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("required field 'title'")
		}
		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("title length must be between 1 and 100, got %d", titleLen)
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("description length must be between 1 and 1000, got %d", descriptionLen)
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("required field 'completed' can't be null")
		}
	}

	return nil
}

// PatchTask godoc
// @Summary Обновить задачу
// @Description Обновление информации о существующей в системе задаче
// @Description ### Логика обновления полей (Three-state logic):
// @Description 1. **Поле не передано**: `description` игнорируется, значение в БД не меняется
// @Description 2. **Явно передано значение**: `"description":"какое то описание"` - устанавливается новое значение description  в БД
// @Description 3. **Передан null**: `"description" : null` - очищает поле в БД (set to NULL)
// @Description Ограничения `title` и `completed` не могут быть выставлены в null
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID изменяемой задачи"
// @Param request body PatchTaskRequest true "PatchTask тело запроса"
// @Success 200 {object} PatchTaskResponse "Успешно измененная задача"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err, "failed to get task id")
		return
	}

	var req PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err, "failed to decode and validate http request")
		return
	}

	taskPatch := TaskPatchFromRequest(req)
	taskDomain, err := h.tasksService.PatchTask(ctx, id, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err, "failed to patch task")
		return
	}

	response := PatchTaskResponse(TaskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func TaskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
