package controllers

import (
	"app/models"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type createProduct struct {
	Name  string                `form:"name" validate:"required"`
	Desc  string                `form:"desc" validate:"required"`
	Price int                   `form:"price" validate:"required"`
	Image *multipart.FileHeader `form:"image" validate:"required"`
}

type productResponse struct {
	Name  string                `form:"name" validate:"required"`
	Desc  string                `form:"desc" validate:"required"`
	Price int                   `form:"price" validate:"required"`
	Image *multipart.FileHeader `form:"image" validate:"required"`
}

type Product struct {
	DB *gorm.DB
}

func (p *Product) FindAll(ctx *fiber.Ctx) error {
	// db := config.GetDB()
	// name := ctx.Params("name")
	// msg := fmt.Sprintf("Hello, %s", name)

	var products []models.Product

	p.DB.Find(&products)
	return ctx.JSON(map[string]interface{}{"message": products})
}

func (p *Product) Create(ctx *fiber.Ctx) error {
	// db := config.GetDB()
	var form createProduct
	if err := ctx.BodyParser(&form); err != nil {
		ctx.Status(fiber.StatusUnprocessableEntity).JSON(map[string]string{"error": err.Error()})
		return err
	}

	var product models.Product

	product.Name = form.Name
	product.Desc = form.Desc
	product.Price = form.Price

	p.DB.Create(&product)

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{"product": product})
}
