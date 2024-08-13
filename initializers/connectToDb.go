package initializers

import (
	"log"
	"main-module/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
	dsn := os.Getenv("DB")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
}