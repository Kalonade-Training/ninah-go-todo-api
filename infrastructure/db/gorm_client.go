package db

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MustOpen() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB open failed: %v", err)
	}
	if err := db.AutoMigrate(&TodoModel{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	return db
}
