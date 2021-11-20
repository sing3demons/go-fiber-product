package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sing3demons/product-app/config"
	"github.com/sing3demons/product-app/routes"
	"github.com/sing3demons/product-app/seeds"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDB()
	seeds.Load()

	app := fiber.New()
	routes.Serve(app)
	app.Static("/uploads", "./uploads")

	//สร้าง folder
	uploadDirs := [...]string{"products", "users"}
	for _, dir := range uploadDirs {
		os.MkdirAll("uploads/"+dir, 0755)
	}

	app.Listen(":" + os.Getenv("PORT"))
}
