package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/ninahf618/go-todo-api/infrastructure/db"
	"github.com/ninahf618/go-todo-api/interfaces/handler"
	"github.com/ninahf618/go-todo-api/middleware"
	"github.com/ninahf618/go-todo-api/pkg/auth"
	"github.com/ninahf618/go-todo-api/pkg/security"
	"github.com/ninahf618/go-todo-api/usecases"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8091"
	}
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set (put it in .env or container env)")
	}

	gormDB := db.InitDB()
	if gormDB == nil {
		log.Fatal("DB init failed: gormDB is nil")
	}

	userRepo := db.NewUserRepository(gormDB)
	todoRepo := db.NewTodoRepository(gormDB)

	tokenSvc := auth.NewJWTService(os.Getenv("JWT_SECRET"))
	pwdSvc := security.NewBcryptService()

	userUC := usecases.NewUserUsecase(userRepo, tokenSvc, pwdSvc)
	todoUC := usecases.NewTodoUsecase(todoRepo)

	userH := handler.NewUserHandler(userUC)
	todoH := handler.NewTodoHandler(todoUC)

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Go TodoAPI is running"})
	})

	r.POST("/register", userH.Register)
	r.POST("/login", userH.Login)

	authDbg := r.Group("/auth")
	authDbg.Use(middleware.JWTMiddleware())
	authDbg.GET("/whoami", func(c *gin.Context) {
		if v, ok := c.Get("user_id"); ok {
			if s, ok := v.(string); ok {
				c.JSON(http.StatusOK, gin.H{"user_id": s})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"user_id": nil})
	})

	authGroup := r.Group("/todos")
	authGroup.Use(middleware.JWTMiddleware()) // no args
	{
		authGroup.GET("", todoH.List)
		authGroup.GET("/:id", todoH.Detail)
		authGroup.POST("", todoH.Create)
		authGroup.PATCH("/:id", todoH.Update)
		authGroup.DELETE("/:id", todoH.Delete)
		authGroup.POST("/:id/duplicate", todoH.Duplicate)
	}

	log.Printf("Server is running on port %s", port)
	log.Fatal(r.Run(":" + port))
}
