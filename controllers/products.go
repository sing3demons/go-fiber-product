package controllers

import (
	"app/config"
	"app/models"
	"fmt"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type productForm struct {
	Name  string                `form:"name" validate:"required"`
	Desc  string                `form:"desc" validate:"required"`
	Price int                   `form:"price" validate:"required"`
	Image *multipart.FileHeader `form:"image" validate:"required"`
}

type updateProductForm struct {
	Name  string                `form:"name"`
	Desc  string                `form:"desc"`
	Price int                   `form:"price"`
	Image *multipart.FileHeader `form:"image"`
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

	p.DB.Order("id desc").Find(&products)
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
	var form productForm
	if err := ctx.BodyParser(&form); err != nil {
		ctx.Status(fiber.StatusUnprocessableEntity).JSON(config.H{"error": err.Error()})
		return err
	}

	var product models.Product
	copier.Copy(&product, &form)
	p.DB.Create(&product)

	p.setProductImage(ctx, &product)

	var serializedProduct productRespons
	copier.Copy(&serializedProduct, &product)

	return ctx.Status(fiber.StatusOK).JSON(config.H{"product": serializedProduct})
}

func (p *Product) Update(ctx *fiber.Ctx) error {
	var form updateProductForm
	if err := ctx.BodyParser(&form); err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(config.H{"error": err.Error()})
	}

	product, err := p.findProductByID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(config.H{"error": err.Error()})
	}

	copier.Copy(&product, &form)

	if err := p.DB.Save(&product).Error; err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(config.H{"error": err.Error()})
	}

	p.setProductImage(ctx, product)

	return ctx.SendStatus(fiber.StatusOK)
}

func (p *Product) findProductByID(ctx *fiber.Ctx) (*models.Product, error) {
	var product models.Product
	id := ctx.Params("id")

	if err := p.DB.First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil

}

func (p *Product) setProductImage(ctx *fiber.Ctx, product *models.Product) error {
	file, err := ctx.FormFile("image")
	if err != nil || file == nil {
		return err
	}

	if product.Image != "" {
		product.Image = strings.Replace(product.Image, os.Getenv("HOST"), "", 1)
		fmt.Println(product.Image)
		pwd, _ := os.Getwd()
		fmt.Println(pwd)
		os.Remove(pwd + product.Image)

	}

	path := "uploads/products/" + strconv.Itoa(int(product.ID))
	os.MkdirAll(path, 0755)

	filename := path + "/" + file.Filename
	if err := ctx.SaveFile(file, filename); err != nil {
		return err
	}

	product.Image = os.Getenv("HOST") + "/" + filename

	if err := p.DB.Save(product).Error; err != nil {
		return err
	}

	return nil

}
