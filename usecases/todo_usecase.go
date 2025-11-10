package usecases

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ninahf618/go-todo-api/domain/entities"
	"github.com/ninahf618/go-todo-api/domain/repositories"
	"github.com/ninahf618/go-todo-api/domain/valueobjects"
)

type TodoUsecase interface {
	List(userID string, f repositories.TodoListFilter) ([]entities.Todo, error)
	Detail(userID, id string) (*entities.Todo, error)
	Create(userID, title, body string, due *time.Time) (*entities.Todo, error)
	Update(userID, id string, title *string, body *string, due *time.Time, completed *bool) (*entities.Todo, error)
	Delete(userID, id string) error
	Duplicate(userID, id string) (*entities.Todo, error)
}

type todoUsecase struct{ repo repositories.TodoRepository }

func NewTodoUsecase(r repositories.TodoRepository) TodoUsecase {
	return &todoUsecase{repo: r}
}

func (u *todoUsecase) List(userID string, f repositories.TodoListFilter) ([]entities.Todo, error) {
	f.UserID = userID
	return u.repo.ListByFilter(f)
}

func (u *todoUsecase) Detail(userID, id string) (*entities.Todo, error) {
	return u.repo.FindByID(id, userID)
}

func (u *todoUsecase) Create(userID, title, body string, due *time.Time) (*entities.Todo, error) {
	name, err := valueobjects.NewName(title)
	if err != nil {
		return nil, err
	}
	b, err := valueobjects.NewBody(body)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	ent := &entities.Todo{
		ID:        uuid.NewString(),
		UserID:    userID,
		Title:     name,
		Body:      b,
		DueDate:   due,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := u.repo.Create(ent); err != nil {
		return nil, err
	}
	return ent, nil
}

func (u *todoUsecase) Update(userID, id string, title *string, body *string, due *time.Time, completed *bool) (*entities.Todo, error) {
	ent, err := u.repo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}
	if ent == nil {
		return nil, errors.New("todo not found")
	}

	if title != nil {
		if v, err := valueobjects.NewName(*title); err != nil {
			return nil, err
		} else {
			ent.Title = v
		}
	}
	if body != nil {
		if v, err := valueobjects.NewBody(*body); err != nil {
			return nil, err
		} else {
			ent.Body = v
		}
	}
	if due != nil {
		ent.DueDate = due
	}
	if completed != nil {
		ent.Completed = *completed
		if ent.Completed {
			now := time.Now().UTC()
			ent.CompletedAt = &now
		} else {
			ent.CompletedAt = nil
		}
	}
	ent.UpdatedAt = time.Now().UTC()

	if err := u.repo.Update(ent); err != nil {
		return nil, err
	}
	return ent, nil
}

func (u *todoUsecase) Delete(userID, id string) error {
	return u.repo.Delete(id, userID)
}

func (u *todoUsecase) Duplicate(userID, id string) (*entities.Todo, error) {
	src, err := u.repo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}
	if src == nil {
		return nil, errors.New("todo not found")
	}

	newTitle := src.Title.String() + "のコピー"
	if len([]rune(newTitle)) > 50 {
		return nil, errors.New("Todo title too long, please shorten the original todo")
	}

	newID := uuid.NewString()
	return u.repo.Duplicate(id, userID, newID, newTitle)
}
