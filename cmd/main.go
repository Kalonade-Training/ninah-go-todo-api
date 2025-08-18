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
	"github.com/ninahf618/go-todo-api/usecases"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, proceeding with system environment variables")
	}

	log.Printf("DATABASE_URL = %q", os.Getenv("DATABASE_URL"))
	gormDB := db.MustOpen()

	r := gin.Default()

	userRepo := db.NewUserRepository(gormDB)
	todoRepo := db.NewTodoRepository(gormDB)

	userUC := usecases.NewUserUsecase(userRepo)
	todoUC := usecases.NewTodoUsecase(todoRepo)

	userHandler := handler.NewUserHandler(userUC)
	todoHandler := handler.NewTodoHandler(todoUC)

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Go TodoAPI is running"})
	})

	auth := r.Group("/todos")
	auth.Use(middleware.JWTMiddleware())
	{
		auth.GET("", todoHandler.List)
		auth.POST("", todoHandler.Create)
		auth.PUT("/:id", todoHandler.Update)
		auth.DELETE("/:id", todoHandler.Delete)
	}

	if gormDB == nil {
		log.Fatal("global DB is nil")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)

	log.Fatal(r.Run(":" + port))

}
