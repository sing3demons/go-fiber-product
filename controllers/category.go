package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/product-app/models"
	"gorm.io/gorm"
)

type Category struct {
	DB *gorm.DB
}

type categoryForm struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type categoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type categoryProductResponse struct {
	Name    string `json:"name"`
	Product []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"products,omitempty"`
}

func (c *Category) FindAll(ctx *fiber.Ctx) error {
	var categories []models.Category

	if err := c.DB.Find(&categories).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	if s := ctx.Query("name"); s != "" {
		c.DB.Where("name = ?", s).Find(&categories)
	}

	serializedCategory := []categoryResponse{}
	copier.Copy(&serializedCategory, &categories)

	resp := fiber.Map{
		"category": serializedCategory,
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (c *Category) FindCategoryProduct(ctx *fiber.Ctx) error {
	var categories []models.Category

	if err := c.DB.Preload("Product", func(db *gorm.DB) *gorm.DB {
		var products []models.Product
		page, _ := strconv.Atoi(ctx.Query("page", "1"))
		limit, _ := strconv.Atoi(ctx.Query("limit", "24"))

		offset := (page - 1) * limit

		return db.Limit(limit).Offset(offset).Order("id desc").Find(&products)
	}).Find(&categories).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	if s := ctx.Query("name"); s != "" {
		c.DB.Where("name = ?", s).Find(&categories)
	}

	serializedCategory := []categoryProductResponse{}
	copier.Copy(&serializedCategory, &categories)

	resp := fiber.Map{
		"category": serializedCategory,
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (c *Category) FindOne(ctx *fiber.Ctx) error {
	category, err := c.findByCategoryID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	serializedCategory := categoryResponse{}
	copier.Copy(&serializedCategory, &category)
	resp := fiber.Map{
		"category": serializedCategory,
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (c *Category) Create(ctx *fiber.Ctx) error {
	var form categoryForm
	if err := ctx.BodyParser(&form); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var category models.Category
	copier.Copy(&category, &form)

	if err := c.DB.Create(&category).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.SendStatus(fiber.StatusCreated)
}

func (c *Category) Update(ctx *fiber.Ctx) error {
	var form categoryForm
	if err := ctx.BodyParser(&form); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	category, err := c.findByCategoryID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	copier.Copy(&category, &form)
	c.DB.Model(category).Update("name", &category.Name)
	return ctx.SendStatus(fiber.StatusNoContent)
}
func (c *Category) Delete(ctx *fiber.Ctx) error {
	category, err := c.findByCategoryID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	c.DB.Delete(category)
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *Category) findByCategoryID(ctx *fiber.Ctx) (*models.Category, error) {
	var category models.Category
	id, _ := ctx.ParamsInt("id")
	if err := c.DB.Preload("Product").First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
