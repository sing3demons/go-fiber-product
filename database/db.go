package database

import (
	"fmt"
	"os"

	"github.com/sing3demons/product-app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s  sslmode=disable TimeZone=Asia/Bangkok", dbHost, user, password, dbName, db_port)
	database, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic("failed to connect database")
	}

	database.AutoMigrate(&models.Product{}, &models.Category{})

	db = database
}

func GetDB() *gorm.DB {
	return db
}
