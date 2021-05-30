package config

import (
	"app/models"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	// database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	HOST := os.Getenv("DB_HOST")
	DBNAME := os.Getenv("DB_NAME")

	URL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS,
		HOST, DBNAME)

	database, err := gorm.Open(mysql.Open(URL))

	if err != nil {
		panic("failed to connect database")
	}

	database.AutoMigrate(&models.Product{})

	db = database
}

func GetDB() *gorm.DB {
	return db
}
