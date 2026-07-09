package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IvanJSBog/goland-todo-app/internal/core/domain"
	core_logger "github.com/IvanJSBog/goland-todo-app/internal/core/logger"
	core_http_request "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created" example:"10"`
	TasksCompleted             int      `json:"tasks_completed" example:"12"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate" example:"50"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time" example:"1m30s"`
}

// GetStatistics godoc
// @Summary Получение статистики
// @Description Получение статистики по задачам с опциональной фильтрацией по user_id и/или по временному промежутку
// @Tags statistics
// @Produce json
// @Param user_id query int false "Фильтрация статистики по конкретному пользователю"
// @Param from query string false "Начало временного промежутка, формат YYYY-MM-DD"
// @Param to query string false "Конец временного промежутка, формат YYYY-MM-DD"
// @Success 200 {object} GetStatisticsResponse "Успешное получение статистики"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /statistics [get]
func (h *StatisticsTransportHTTP) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, rw)

	userID, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user_id/from/to query params")
		return
	}

	statisticsDomain, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")
		return
	}

	response := toDTOFromDomain(statisticsDomain)
	responseHandler.JSONResponse(response, http.StatusOK)

}

func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	userID, err := core_http_request.GetIntQueryParam(r, "user_id")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get user_id query param: %w", err)
	}

	from, err := core_http_request.GetDateQueryParam(r, "from")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get from query param: %w", err)
	}

	to, err := core_http_request.GetDateQueryParam(r, "to")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get to query param: %w", err)
	}

	return userID, from, to, nil
}

func toDTOFromDomain(statisticsDomain domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statisticsDomain.TasksAverageCompletionTime != nil {
		duration := statisticsDomain.TasksAverageCompletionTime.String()
		avgTime = &duration
	}
	return GetStatisticsResponse{
		TasksCreated:               statisticsDomain.TasksCreated,
		TasksCompleted:             statisticsDomain.TasksCompleted,
		TasksCompletedRate:         statisticsDomain.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}
