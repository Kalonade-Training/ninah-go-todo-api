package usecases

import (
	"github.com/google/uuid"
	"github.com/ninahf618/go-todo-api/domain/entities"
	"github.com/ninahf618/go-todo-api/domain/repositories"
)

type TodoUsecase struct {
	repo repositories.TodoRepository
}

func NewTodoUsecase(repo repositories.TodoRepository) *TodoUsecase {
	return &TodoUsecase{repo: repo}
}

func (u *TodoUsecase) Create(title, description, userID string) (*entities.Todo, error) {
	todo := &entities.Todo{
		ID:          uuid.NewString(),
		Title:       title,
		Description: description,
		UserID:      userID,
	}
	err := u.repo.Create(todo)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (u *TodoUsecase) ListByUserID(userID string) ([]entities.Todo, error) {
	return u.repo.FindAllByUserID(userID)
}

func (u *TodoUsecase) Update(id, title, description string, completed bool) error {
	todo, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}
	todo.Title = title
	todo.Description = description
	todo.Completed = completed
	return u.repo.Update(todo)
}

func (u *TodoUsecase) Delete(id string) error {
	return u.repo.Delete(id)
}

func (u *TodoUsecase) GetByID(id string) (*entities.Todo, error) {
	return u.repo.FindByID(id)
}
