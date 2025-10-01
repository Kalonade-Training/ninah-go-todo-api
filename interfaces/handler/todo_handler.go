package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ninahf618/go-todo-api/domain/repositories"
	"github.com/ninahf618/go-todo-api/usecases"
)

type TodoHandler struct{ uc usecases.TodoUsecase }

func NewTodoHandler(uc usecases.TodoUsecase) *TodoHandler { return &TodoHandler{uc: uc} }

func userIDFrom(c *gin.Context) string {
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	if v, ok := c.Get("userID"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (h *TodoHandler) List(c *gin.Context) {
	uid := userIDFrom(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var f repositories.TodoListFilter
	f.TitleLike = c.Query("title")
	f.BodyLike = c.Query("body")

	parseDate := func(s string) *time.Time {
		if s == "" {
			return nil
		}
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return nil
		}
		return &t
	}
	f.DueFrom = parseDate(c.Query("due_from"))
	f.DueTo = parseDate(c.Query("due_to"))

	if v := c.Query("completed"); v != "" {
		if v == "true" || v == "1" {
			t := true
			f.Completed = &t
		} else if v == "false" || v == "0" {
			fa := false
			f.Completed = &fa
		}
	}

	items, err := h.uc.List(uid, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *TodoHandler) Detail(c *gin.Context) {
	uid := userIDFrom(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	t, err := h.uc.Detail(uid, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if t == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *TodoHandler) Create(c *gin.Context) {
	uid := userIDFrom(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var body struct {
		Title   string `json:"title" binding:"required"`
		Body    string `json:"body"`
		DueDate string `json:"due_date"` // YYYY-MM-DD (optional)
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "BadRequest", "message": err.Error()})
		return
	}
	var due *time.Time
	if body.DueDate != "" {
		if t, err := time.Parse("2006-01-02", body.DueDate); err == nil {
			due = &t
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "BadRequest", "message": "invalid due_date format (use YYYY-MM-DD)"})
			return
		}
	}
	ent, err := h.uc.Create(uid, body.Title, body.Body, due)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ent)
}

func (h *TodoHandler) Update(c *gin.Context) {
	uid := userIDFrom(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	var body struct {
		Title     *string `json:"title"`
		Body      *string `json:"body"`
		DueDate   *string `json:"due_date"`  // YYYY-MM-DD
		Completed *bool   `json:"completed"` // true/false
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "BadRequest", "message": err.Error()})
		return
	}
	var due *time.Time
	if body.DueDate != nil {
		if *body.DueDate == "" {
			due = nil
		} else {
			t, err := time.Parse("2006-01-02", *body.DueDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "BadRequest", "message": "invalid due_date format (use YYYY-MM-DD)"})
				return
			}
			due = &t
		}
	}
	ent, err := h.uc.Update(uid, id, body.Title, body.Body, due, body.Completed)
	if err != nil {
		if err.Error() == "todo not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ent)
}

func (h *TodoHandler) Delete(c *gin.Context) {
	uid := userIDFrom(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	if err := h.uc.Delete(uid, id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *TodoHandler) Duplicate(c *gin.Context) {
	uid := userIDFrom(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	ent, err := h.uc.Duplicate(uid, id)
	if err != nil {
		if err.Error() == "todo not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ent)
}
