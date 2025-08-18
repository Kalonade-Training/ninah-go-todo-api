package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ninahf618/go-todo-api/domain/entities"
	"github.com/ninahf618/go-todo-api/usecases"
)

type TodoHandler struct{ uc *usecases.TodoUsecase }

func NewTodoHandler(uc *usecases.TodoUsecase) *TodoHandler {
	return &TodoHandler{uc: uc}
}

type createBody struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

func (h *TodoHandler) Create(c *gin.Context) {
	var b createBody
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "BadRequest", "message": err.Error()})
		return
	}

	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	desc := ""
	if b.Description != nil {
		desc = *b.Description
	}

	created, err := h.uc.Create(b.Title, desc, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ValidationError", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": toJSON(created)})
}

type updateBody struct {
	Title       *string  `json:"title"`
	Description **string `json:"description"`
	Completed   *bool    `json:"completed"`
}

func (h *TodoHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var b updateBody
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "BadRequest", "message": err.Error()})
		return
	}

	updated, err := h.uc.Update(id, b.Title, b.Description, b.Completed)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ValidationError", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toJSON(updated)})
}

func (h *TodoHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.uc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal", "message": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TodoHandler) List(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	limit := atoiDefault(c.Query("limit"), 20, 1, 100)
	offset := atoiDefault(c.Query("offset"), 0, 0, 1<<31-1)
	q := c.Query("q")
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	rows, total, err := h.uc.ListByUserID(userID, limit, offset, q, sort, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal", "message": err.Error()})
		return
	}

	out := make([]any, 0, len(rows))
	for _, e := range rows {
		out = append(out, toJSON(e))
	}
	c.JSON(http.StatusOK, gin.H{
		"data": out,
		"meta": gin.H{"total": total, "limit": limit, "offset": offset},
	})
}

func toJSON(t *entities.Todo) map[string]any {
	var desc *string
	if t.Description() != nil {
		s := t.Description().String()
		desc = &s
	}
	return map[string]any{
		"id":          t.ID().String(),
		"title":       t.Title().String(),
		"description": desc,
		"userId":      t.UserID().String(),
		"completed":   t.Completed(),
		"createdAt":   t.CreatedAt(),
		"updatedAt":   t.UpdatedAt(),
	}
}

func atoiDefault(s string, def, min, max int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

func userIDFromContext(c *gin.Context) (string, bool) {
	if v, ok := c.Get("userID"); ok {
		switch vv := v.(type) {
		case string:
			if _, err := uuid.Parse(vv); err == nil {
				return vv, true
			}
		case uuid.UUID:
			return vv.String(), true
		}
	}

	if v, ok := c.Get("user_id"); ok {
		if s, ok2 := v.(string); ok2 {
			if _, err := uuid.Parse(s); err == nil {
				return s, true
			}
		}
	}
	return "", false
}
