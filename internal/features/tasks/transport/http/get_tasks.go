package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/IvanJSBog/goland-todo-app/internal/core/logger"
	core_http_request "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, rw)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId limit offset query params")
		return
	}

	tasksDomains, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks")
		return
	}

	response := GetTasksResponse(TaskDTOsFromDomains(tasksDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	userID, err := core_http_request.GetIntQueryParam(r, "user_id")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get user id parameter: %w", err)
	}
	limit, err := core_http_request.GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get limit parameter: %w", err)
	}
	offset, err := core_http_request.GetIntQueryParam(r, "offset")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get offset parameter: %w", err)
	}

	return userID, limit, offset, nil
}
