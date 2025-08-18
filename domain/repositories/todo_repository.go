package repositories

import "github.com/ninahf618/go-todo-api/domain/entities"

type ListParams struct {
	Limit  int
	Offset int
	Sort   string
	Order  string
	Q      string
	UserID string
}

type TodoRepository interface {
	FindAll(q ListParams) (rows []*entities.Todo, total int64, err error)
	FindByID(id string) (*entities.Todo, error)
	Create(t *entities.Todo) (*entities.Todo, error)
	Update(t *entities.Todo) (*entities.Todo, error)
	Delete(id string) error
}
