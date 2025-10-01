package entities

import (
	"time"

	"github.com/ninahf618/go-todo-api/domain/valueobjects"
)

type Todo struct {
	ID          string
	UserID      string
	Title       valueobjects.Name
	Body        valueobjects.Body
	DueDate     *time.Time
	Completed   bool
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
