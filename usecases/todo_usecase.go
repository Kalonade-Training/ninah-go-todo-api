package usecases

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ninahf618/go-todo-api/domain/entities"
	"github.com/ninahf618/go-todo-api/domain/repositories"
	"github.com/ninahf618/go-todo-api/domain/vo"
)

type TodoUsecase struct {
	repo repositories.TodoRepository
}

func NewTodoUsecase(repo repositories.TodoRepository) *TodoUsecase {
	return &TodoUsecase{repo: repo}
}

func (u *TodoUsecase) Create(title, description, userID string) (*entities.Todo, error) {
	titleVO, err := vo.NewTitle(title)
	if err != nil {
		return nil, err
	}

	var descVO *vo.Description
	if description != "" {
		d, err := vo.NewDescription(description)
		if err != nil {
			return nil, err
		}
		descVO = &d
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID: %w", err)
	}

	e := entities.NewTodo(titleVO, uid, descVO, nil)
	return u.repo.Create(e)
}

func (u *TodoUsecase) ListByUserID(userID string, limit, offset int, q, sort, order string) ([]*entities.Todo, int64, error) {
	params := repositories.ListParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
		Q:      q,
		Sort:   sort,
		Order:  order,
	}
	return u.repo.FindAll(params)
}

func (u *TodoUsecase) Update(id string, title *string, description **string, completed *bool) (*entities.Todo, error) {
	t, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	var upd entities.TodoUpdate

	if title != nil {
		tv, err := vo.NewTitle(*title)
		if err != nil {
			return nil, err
		}
		upd.Title = &tv
	}

	if description != nil {
		if *description == nil {
			var nilDesc *vo.Description = nil
			upd.Description = &nilDesc
		} else {
			dv, err := vo.NewDescription(**description)
			if err != nil {
				return nil, err
			}
			descPtr := &dv
			upd.Description = &descPtr
		}
	}

	if completed != nil {
		upd.Completed = completed
	}

	t.UpdateValues(upd)
	return u.repo.Update(t)
}

func (u *TodoUsecase) Delete(id string) error {
	return u.repo.Delete(id)
}

func (u *TodoUsecase) GetByID(id string) (*entities.Todo, error) {
	return u.repo.FindByID(id)
}
