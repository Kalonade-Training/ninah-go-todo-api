package db

import "time"

type TodoModel struct {
	ID          string     `gorm:"type:char(36);primaryKey"`
	UserID      string     `gorm:"type:char(36);index;not null"`
	Title       string     `gorm:"size:50;not null;index"`
	Body        string     `gorm:"size:1000;index"`
	DueDate     *time.Time `gorm:"index"`
	Completed   bool       `gorm:"not null;default:false;index"`
	CompletedAt *time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime;index"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (TodoModel) TableName() string { return "todos" }
