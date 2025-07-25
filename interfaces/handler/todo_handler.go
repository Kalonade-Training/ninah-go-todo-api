package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ninahf618/go-todo-api/usecases"
)

type TodoHandler struct {
	usecase *usecases.TodoUsecase
}

func NewTodoHandler(usecase *usecases.TodoUsecase) *TodoHandler {
	return &TodoHandler{usecase}
}

func (h *TodoHandler) List(c *gin.Context) {
	userID := c.GetString("userID")

	todos, err := h.usecase.ListByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")

	var body struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(body.Title) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title cannot be empty"})
		return
	}
	if len(body.Title) > 255 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is too long"})
		return
	}
	if len(body.Description) > 2000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Description is too long"})
		return
	}

	todo, err := h.usecase.Create(body.Title, body.Description, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) Update(c *gin.Context) {
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(body.Title) != "" && len(body.Title) > 255 {
		if len(body.Title) > 255 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title is too long"})
			return
		}
	}
	if len(body.Description) > 2000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Description is too long"})
		return
	}

	id := c.Param("id")

	err := h.usecase.Update(id, body.Title, body.Description, body.Completed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated"})
}

func (h *TodoHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")

	todo, err := h.usecase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	if todo.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	err = h.usecase.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})

}
