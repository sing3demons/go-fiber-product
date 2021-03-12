package routes

import (
	"app/config"
	"app/controllers"

	"github.com/gofiber/fiber/v2"
)

func Serve() *fiber.App {
	db := config.GetDB()
	app := fiber.New()
	app.Get("", homepage)
	v1 := app.Group("api/v1")

	productController := controllers.Product{DB: db}
	productsGroup := v1.Group("/products")
	{
		productsGroup.Get("", productController.FindAll)
		productsGroup.Post("", productController.Create)
	}

	return app
}

// http://127.0.0.1:8080
func homepage(ctx *fiber.Ctx) error {
	name := ctx.Query("name")

	if name == "" {
		name = ", world"
	}

	return ctx.SendString("Hello " + name)
}
