package tasks_postgres_repository

import (
	"time"

	"github.com/IvanJSBog/goland-todo-app/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func TaskDomainsFromModels(taskModels []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(taskModels))
	for i, t := range taskModels {
		domains[i] = domain.NewTask(t.ID, t.Version, t.Title, t.Description, t.Completed, t.CreatedAt, t.CompletedAt, t.AuthorUserID)
	}
	return domains
}
