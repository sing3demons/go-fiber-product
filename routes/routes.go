package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/product-app/config"
	"github.com/sing3demons/product-app/controllers"
)

func Serve(app *fiber.App) {
	db := config.GetDB()

	app.Get("", homepage)
	v1 := app.Group("api/v1")

	productController := controllers.Product{DB: db}
	productsGroup := v1.Group("/products")
	{
		productsGroup.Get("", productController.FindAll)
		productsGroup.Get("/:id", productController.FindOne)
		productsGroup.Post("", productController.Create)
		productsGroup.Put("/:id", productController.Update)
		productsGroup.Delete("/:id", productController.Delete)
	}

	categoryController := controllers.Category{DB: db}
	categoryGroup := v1.Group("/categories")
	categoryGroup.Get("", categoryController.FindAll)
	categoryGroup.Get("/products", categoryController.FindCategoryProduct)
	categoryGroup.Get("/:id", categoryController.FindOne)
	// categoryGroup.Use(authenticate, authorize)
	{
		categoryGroup.Post("", categoryController.Create)
		categoryGroup.Patch("/:id", categoryController.Update)
		categoryGroup.Delete("/:id", categoryController.Delete)
	}

}

// http://127.0.0.1:8080
func homepage(ctx *fiber.Ctx) error {
	name := ctx.Query("name")

	if name == "" {
		name = ", world"
	}

	return ctx.SendString("Hello " + name)
}
