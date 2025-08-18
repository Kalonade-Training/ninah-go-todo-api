package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/ninahf618/go-todo-api/domain/vo"
)

type TodoProps struct {
	ID          uuid.UUID
	Title       vo.Title
	Description *vo.Description
	Completed   bool
	UserID      uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Todo struct {
	id          uuid.UUID
	title       vo.Title
	description *vo.Description
	userID      uuid.UUID
	completed   bool
	createdAt   time.Time
	updatedAt   time.Time
}

func NewTodo(title vo.Title, userID uuid.UUID, desc *vo.Description, due *time.Time) *Todo {
	now := time.Now().UTC()
	return &Todo{
		id:          uuid.New(),
		title:       title,
		description: desc,
		userID:      userID,
		completed:   false,
		createdAt:   now,
		updatedAt:   now,
	}
}

func RebuildTodo(p TodoProps) *Todo {
	return &Todo{
		id:          p.ID,
		title:       p.Title,
		description: p.Description,
		userID:      p.UserID,
		completed:   p.Completed,
		createdAt:   p.CreatedAt,
		updatedAt:   p.UpdatedAt,
	}
}

func (t *Todo) ID() uuid.UUID                { return t.id }
func (t *Todo) Title() vo.Title              { return t.title }
func (t *Todo) Description() *vo.Description { return t.description }
func (t *Todo) UserID() uuid.UUID            { return t.userID }
func (t *Todo) Completed() bool              { return t.completed }
func (t *Todo) CreatedAt() time.Time         { return t.createdAt }
func (t *Todo) UpdatedAt() time.Time         { return t.updatedAt }

type TodoUpdate struct {
	Title       *vo.Title
	Description *(*vo.Description)
	Completed   *bool
}

func (t *Todo) UpdateValues(u TodoUpdate) {
	if u.Title != nil {
		t.title = *u.Title
	}
	if u.Description != nil {
		t.description = *u.Description
	}
	if u.Completed != nil {
		t.completed = *u.Completed
	}
	t.updatedAt = time.Now().UTC()
}
