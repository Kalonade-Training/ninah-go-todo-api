package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
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

	gormDB := db.InitDB()

	userRepo := db.NewUserRepository(gormDB)
	todoRepo := db.NewTodoRepository(gormDB)

	tokenSvc := auth.NewJWTService(os.Getenv("JWT_SECRET"))
	pwdSvc := security.NewBcryptService()

	userUC := usecases.NewUserUsecase(userRepo, tokenSvc, pwdSvc)
	todoUC := usecases.NewTodoUsecase(todoRepo)

	userH := handler.NewUserHandler(userUC)
	todoH := handler.NewTodoHandler(todoUC)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:5174",
			"http://localhost:5175",
			"https://ninah-todo-frontend-new.vercel.app",
			"https://ninah-todo-frontend-a4sfilkps-ninahs-projects-5aeacde4.vercel.app",
			"https://ninah-todo-frontend-new-git-rew-9634db-ninahs-projects-5aeacde4.vercel.app",
			"https://*.vercel.app",
		},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	r.POST("/register", userH.Register)
	r.POST("/login", userH.Login)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Go TodoAPI is running"})
	})

	authGroup := r.Group("/todos")
	authGroup.Use(middleware.JWTMiddleware())
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
