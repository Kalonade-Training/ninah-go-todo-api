package db

import (
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ninahf618/go-todo-api/domain/entities"
	"github.com/ninahf618/go-todo-api/domain/repositories"
	"github.com/ninahf618/go-todo-api/domain/vo"
	"gorm.io/gorm"
)

type TodoModel struct {
	ID          string    `gorm:"primaryKey;type:char(36)"`
	Title       string    `gorm:"size:200;not null"`
	Description *string   `gorm:"type:text"`
	UserID      string    `gorm:"type:char(36);index"`
	Completed   bool      `gorm:"not null;default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

var _ repositories.TodoRepository = (*GormTodoRepository)(nil)

func toEntity(m *TodoModel) (*entities.Todo, error) {
	title, err := vo.NewTitle(m.Title)
	if err != nil {
		return nil, err
	}
	var desc *vo.Description
	if m.Description != nil {
		d, err := vo.NewDescription(*m.Description)
		if err != nil {
			return nil, err
		}
		desc = &d
	}
	return entities.RebuildTodo(entities.TodoProps{
		ID:          uuid.MustParse(m.ID),
		Title:       title,
		Description: desc,
		UserID:      uuid.MustParse(m.UserID),
		Completed:   m.Completed,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}), nil
}

func fromEntity(e *entities.Todo) *TodoModel {
	var desc *string
	if e.Description() != nil {
		s := e.Description().String()
		desc = &s
	}
	return &TodoModel{
		ID:          e.ID().String(),
		Title:       e.Title().String(),
		Description: desc,
		UserID:      e.UserID().String(),
		Completed:   e.Completed(),
		CreatedAt:   e.CreatedAt(),
		UpdatedAt:   e.UpdatedAt(),
	}
}

type GormTodoRepository struct{ db *gorm.DB }

func NewTodoRepository(gormDB *gorm.DB) repositories.TodoRepository {
	if gormDB == nil {
		log.Fatal("NewTodoRepository received a nil *gorm.DB")
	}
	return &GormTodoRepository{db: gormDB}
}

func (r *GormTodoRepository) FindAll(q repositories.ListParams) ([]*entities.Todo, int64, error) {
	limit, offset := q.Limit, q.Offset
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	orderBy := "created_at desc"
	if q.Sort != "" {
		col := map[string]string{
			"title": "title", "created_at": "created_at", "updated_at": "updated_at", "due_date": "due_date",
		}[strings.ToLower(q.Sort)]
		if col != "" {
			dir := "desc"
			if strings.ToLower(q.Order) == "asc" {
				dir = "asc"
			}
			orderBy = col + " " + dir
		}
	}

	tx := r.db.Model(&TodoModel{})
	if q.UserID != "" {
		tx = tx.Where("user_id = ?", q.UserID)
	}
	if q.Q != "" {
		like := "%" + q.Q + "%"
		tx = tx.Where("title LIKE ? OR description LIKE ?", like, like)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []TodoModel
	if err := tx.Order(orderBy).Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	out := make([]*entities.Todo, 0, len(rows))
	for i := range rows {
		e, err := toEntity(&rows[i])
		if err != nil {
			return nil, 0, err
		}
		out = append(out, e)
	}
	return out, total, nil
}

func (r *GormTodoRepository) FindByID(id string) (*entities.Todo, error) {
	var m TodoModel
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toEntity(&m)
}

func (r *GormTodoRepository) Create(t *entities.Todo) (*entities.Todo, error) {
	m := fromEntity(t)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return toEntity(m)
}

func (r *GormTodoRepository) Update(t *entities.Todo) (*entities.Todo, error) {
	m := fromEntity(t)
	if err := r.db.Model(&TodoModel{}).Where("id = ?", m.ID).
		Updates(map[string]any{
			"title": m.Title, "description": m.Description,
			"completed": m.Completed, "updated_at": time.Now().UTC(),
		}).Error; err != nil {
		return nil, err
	}
	return r.FindByID(m.ID)
}

func (r *GormTodoRepository) Delete(id string) error {
	return r.db.Delete(&TodoModel{}, "id = ?", id).Error
}
