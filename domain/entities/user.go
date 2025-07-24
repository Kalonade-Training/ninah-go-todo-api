package entities

import (
	"time"
)

type User struct {
	ID        string `gorm:"type:char(36);primaryKey"`
	Username  string `json:"username" gorm:"size:100;not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Todos     []Todo `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
