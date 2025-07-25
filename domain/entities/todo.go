package entities

import (
	"time"
)

type Todo struct {
	ID          string    `gorm:"type:char(36);primaryKey"`
	Title       string    `gorm:"not null"`
	Description string    `json:"description"`
	Completed   bool      `gorm:"default:false"`
	UserID      string    `gorm:"type:char(36);not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
