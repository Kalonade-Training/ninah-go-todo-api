package repositories

import (
	"time"

	"github.com/ninahf618/go-todo-api/domain/entities"
)

type TodoListFilter struct {
	UserID    string
	TitleLike string
	BodyLike  string
	DueFrom   *time.Time
	DueTo     *time.Time
	Completed *bool
}

type TodoRepository interface {
	ListByFilter(f TodoListFilter) ([]entities.Todo, error)
	FindByID(id, userID string) (*entities.Todo, error)
	Create(t *entities.Todo) error
	Update(t *entities.Todo) error
	Delete(id, userID string) error
	Duplicate(id, userID, newID, newTitle string) (*entities.Todo, error)
}
