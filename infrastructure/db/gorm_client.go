package db

import (
	"fmt"
	"log"
	"os"

	"github.com/ninahf618/go-todo-api/domain/entities"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	err = db.AutoMigrate(&entities.User{}, &entities.Todo{})
	if err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}

	DB = db
	fmt.Println("Database Connected")
	return db
}
