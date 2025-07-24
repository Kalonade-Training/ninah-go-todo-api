package db

import (
	"log"

	"github.com/ninahf618/go-todo-api/domain/entities"
	"github.com/ninahf618/go-todo-api/domain/repositories"
	"gorm.io/gorm"
)

type GormTodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(gormDB *gorm.DB) repositories.TodoRepository {
	if gormDB == nil {
		log.Fatal("NewTodoRepository received a nil *gorm.DB")
	}
	return &GormTodoRepository{db: gormDB}
}

func (r *GormTodoRepository) Create(todo *entities.Todo) error {
	return r.db.Create(todo).Error
}

func (r *GormTodoRepository) FindAllByUserID(userID string) ([]entities.Todo, error) {
	var todos []entities.Todo
	err := r.db.Where("user_id = ?", userID).Find(&todos).Error
	return todos, err
}

func (r *GormTodoRepository) FindByID(id string) (*entities.Todo, error) {
	var todo entities.Todo
	err := r.db.First(&todo, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *GormTodoRepository) Update(todo *entities.Todo) error {
	return r.db.Save(todo).Error
}

func (r *GormTodoRepository) Delete(id string) error {
	return r.db.Delete(&entities.Todo{}, "id = ?", id).Error
}
