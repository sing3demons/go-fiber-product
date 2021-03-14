package controllers

import (
	"app/config"
	"app/models"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type createProduct struct {
	Name  string                `form:"name" validate:"required"`
	Desc  string                `form:"desc" validate:"required"`
	Price int                   `form:"price" validate:"required"`
	Image *multipart.FileHeader `form:"image" validate:"required"`
}

type productRespons struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Price int    `json:"price"`
	Image string `json:"image"`
}

type Product struct {
	DB *gorm.DB
}

func (p *Product) FindAll(ctx *fiber.Ctx) error {
	var products []models.Product

	p.DB.Find(&products)
	return ctx.Status(fiber.StatusOK).JSON(config.H{"products": products})
}

func (p *Product) FindOne(ctx *fiber.Ctx) error {
	product, err := p.findProductByID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(config.H{"error": err.Error()})
	}

	serializedProduct := productRespons{}
	copier.Copy(&serializedProduct, &product)
	return ctx.Status(fiber.StatusOK).JSON(config.H{"product": serializedProduct})
}

func (p *Product) Create(ctx *fiber.Ctx) error {
	var form createProduct
	if err := ctx.BodyParser(&form); err != nil {
		ctx.Status(fiber.StatusUnprocessableEntity).JSON(config.H{"error": err.Error()})
		return err
	}

	var product models.Product

	product.Name = form.Name
	product.Desc = form.Desc
	product.Price = form.Price

	p.DB.Create(&product)

	return ctx.Status(fiber.StatusOK).JSON(config.H{"product": product})
}

func (p *Product) findProductByID(ctx *fiber.Ctx) (*models.Product, error) {
	var product models.Product
	id := ctx.Params("id")

	if err := p.DB.First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil

}
