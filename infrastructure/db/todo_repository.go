package db

import (
	"errors"
	"time"

	"github.com/ninahf618/go-todo-api/domain/entities"
	"github.com/ninahf618/go-todo-api/domain/repositories"
	"github.com/ninahf618/go-todo-api/domain/valueobjects"
	"gorm.io/gorm"
)

type GormTodoRepository struct{ db *gorm.DB }

func (r *GormTodoRepository) FindAllByUserID(userID string) ([]entities.Todo, error) {
	return r.ListByFilter(repositories.TodoListFilter{UserID: userID})
}

func NewTodoRepository(gormDB *gorm.DB) repositories.TodoRepository {
	return &GormTodoRepository{db: gormDB}
}

func toTodoModel(t *entities.Todo) *TodoModel {
	var due *time.Time
	if t.DueDate != nil {
		d := *t.DueDate
		due = &d
	}
	var completedAt *time.Time
	if t.CompletedAt != nil {
		ca := *t.CompletedAt
		completedAt = &ca
	}
	return &TodoModel{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title.String(),
		Body:        t.Body.String(),
		DueDate:     due,
		Completed:   t.Completed,
		CompletedAt: completedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func toTodoEntity(m *TodoModel) *entities.Todo {
	title, _ := valueobjects.NewName(m.Title)
	body, _ := valueobjects.NewBody(m.Body)
	return &entities.Todo{
		ID:          m.ID,
		UserID:      m.UserID,
		Title:       title,
		Body:        body,
		DueDate:     m.DueDate,
		Completed:   m.Completed,
		CompletedAt: m.CompletedAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *GormTodoRepository) ListByFilter(f repositories.TodoListFilter) ([]entities.Todo, error) {
	q := r.db.Model(&TodoModel{}).Where("user_id = ?", f.UserID)

	if f.TitleLike != "" {
		q = q.Where("title LIKE ?", "%"+f.TitleLike+"%")
	}
	if f.BodyLike != "" {
		q = q.Where("body LIKE ?", "%"+f.BodyLike+"%")
	}
	if f.DueFrom != nil {
		q = q.Where("due_date >= ?", f.DueFrom)
	}
	if f.DueTo != nil {
		q = q.Where("due_date <= ?", f.DueTo)
	}
	if f.Completed != nil {
		q = q.Where("completed = ?", *f.Completed)
	}

	q = q.Order("created_at DESC")

	var rows []TodoModel
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]entities.Todo, 0, len(rows))
	for _, m := range rows {
		out = append(out, *toTodoEntity(&m))
	}
	return out, nil
}

func (r *GormTodoRepository) FindByID(id, userID string) (*entities.Todo, error) {
	var m TodoModel
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toTodoEntity(&m), nil
}

func (r *GormTodoRepository) Create(t *entities.Todo) error {
	return r.db.Create(toTodoModel(t)).Error
}

func (r *GormTodoRepository) Update(t *entities.Todo) error {
	values := map[string]interface{}{
		"title":        t.Title.String(),
		"body":         t.Body.String(), // <-- Body
		"due_date":     t.DueDate,
		"completed":    t.Completed,
		"completed_at": t.CompletedAt,
		"updated_at":   t.UpdatedAt,
	}
	res := r.db.Model(&TodoModel{}).
		Where("id = ? AND user_id = ?", t.ID, t.UserID).
		Updates(values)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *GormTodoRepository) Delete(id, userID string) error {
	res := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&TodoModel{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *GormTodoRepository) Duplicate(id, userID, newID, newTitle string) (*entities.Todo, error) {
	var src TodoModel
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&src).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	now := time.Now().UTC()
	dest := TodoModel{
		ID:          newID,
		UserID:      userID,
		Title:       newTitle,
		Body:        src.Body,
		DueDate:     nil,
		Completed:   false,
		CompletedAt: nil,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := r.db.Create(&dest).Error; err != nil {
		return nil, err
	}
	return toTodoEntity(&dest), nil
}
