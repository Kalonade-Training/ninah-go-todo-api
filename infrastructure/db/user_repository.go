package db

import (
	"errors"
	"log"

	"github.com/ninahf618/go-todo-api/domain/entities"
	"github.com/ninahf618/go-todo-api/domain/repositories"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(gormDB *gorm.DB) repositories.UserRepository {
	if gormDB == nil {
		log.Fatal("NewTodoRepository recieved a nil *gorm.DB")
	}
	return &GormUserRepository{db: gormDB}
}

func (r *GormUserRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByID(id string) (*entities.User, error) {
	var user entities.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
