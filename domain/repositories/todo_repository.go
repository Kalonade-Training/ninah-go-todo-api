package repositories

import (
	"github.com/ninahf618/go-todo-api/domain/entities"
)

type TodoRepository interface {
	Create(Todo *entities.Todo) error
	FindAllByUserID(userID string) ([]entities.Todo, error)
	FindByID(id string) (*entities.Todo, error)
	Update(todo *entities.Todo) error
	Delete(id string) error
}
