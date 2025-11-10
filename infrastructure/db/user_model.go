package db

import "time"

type UserModel struct {
	ID        string    `gorm:"type:char(36);primaryKey"`
	Username  string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:191;uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (UserModel) TableName() string { return "users" }
