package tasks_transport_http

import (
	"time"

	"github.com/IvanJSBog/goland-todo-app/internal/core/domain"
)

type TaskDTOResponse struct {
	ID           int        `json:"id" example:"1"`
	Version      int        `json:"version" example:"1""`
	Title        string     `json:"title" example:"Task Title"`
	Description  *string    `json:"description" example:"Task Description"`
	Completed    bool       `json:"completed" example:"true"`
	CreatedAt    time.Time  `json:"created_at" example:"2020-01-01T00:00:00Z"`
	CompletedAt  *time.Time `json:"completed_at" example:"2020-01-01T00:00:00Z"`
	AuthorUserID int        `json:"author_user_id" example:"1"`
}

func TaskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func TaskDTOsFromDomains(tasks []domain.Task) []TaskDTOResponse {
	taskDTOs := make([]TaskDTOResponse, len(tasks))
	for i := 0; i < len(tasks); i++ {
		taskDTOs[i] = TaskDTOFromDomain(tasks[i])
	}
	return taskDTOs
}
