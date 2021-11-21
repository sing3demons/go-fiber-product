package controllers

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/product-app/cache"
	"github.com/sing3demons/product-app/models"
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
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Price      int    `json:"price"`
	Image      string `json:"image"`
	CategoryID uint   `json:"categoryId"`
	Category   struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
}

type productPaging struct {
	Items  []productRespons `json:"items"`
	Paging *pagingResult    `json:"paging"`
}

type Product struct {
	DB    *gorm.DB
	Redis *cache.Cacher
}

//FindAll - All Products
func (p *Product) FindAll(ctx *fiber.Ctx) error {
	cacheProduct := "items::all"
	cachePage := "items::page"

	cacheItems, err := p.Redis.Get(cacheProduct)
	cacheItemPage, _ := p.Redis.Get(cachePage)
	if err != nil {
		fmt.Println(err)
	}

	if len(cacheItems) > 0 && len(cacheItemPage) > 0 {
		fmt.Println("Get...")

		var items []productRespons
		var page *pagingResult
		if err := json.Unmarshal([]byte(cacheItems), &items); err != nil {
			fmt.Println(err.Error())
			//json: Unmarshal(non-pointer main.Request)
		}
		if err = json.Unmarshal([]byte(cacheItemPage), &page); err != nil {
			fmt.Println(err.Error())
			//json: Unmarshal(non-pointer main.Request)
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"products": productPaging{
				Items:  items,
				Paging: page,
			},
		})

	}

	var products []models.Product

	query := p.DB.Preload("Category").Order("id desc")

	if category := ctx.Query("category"); category != "" {
		c, _ := strconv.Atoi(category)
		query = query.Where("category_id = ?", c)
	}

	pagination := pagination{
		ctx:     ctx,
		query:   p.DB,
		records: &products,
	}
	paging := pagination.pagingResource()

	serializedProduct := []productRespons{}
	copier.Copy(&serializedProduct, &products)

	timeToExpire := 10 * time.Second // 60 * 5 * time.Second // 5m
	p.Redis.Set(cacheProduct, serializedProduct, timeToExpire)
	p.Redis.Set(cachePage, paging, timeToExpire)

	resp := fiber.Map{
		"items": serializedProduct,
		"page":  paging,
	}
	fmt.Println("Set...")
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"products": resp,
	})
}

//FindOne - first product
func (p *Product) FindOne(ctx *fiber.Ctx) error {
	product, err := p.findProductByID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	serializedProduct := productRespons{}
	copier.Copy(&serializedProduct, &product)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"product": serializedProduct})
}

//Create - insert product
func (p *Product) Create(ctx *fiber.Ctx) error {
	var form productForm
	if err := ctx.BodyParser(&form); err != nil {
		ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		return err
	}

	var product models.Product
	copier.Copy(&product, &form)
	p.DB.Create(&product)

	p.setProductImage(ctx, &product)

	var serializedProduct productRespons
	copier.Copy(&serializedProduct, &product)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"product": serializedProduct})
}

//Update - update product
func (p *Product) Update(ctx *fiber.Ctx) error {
	var form updateProductForm
	if err := ctx.BodyParser(&form); err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	product, err := p.findProductByID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	copier.Copy(&product, &form)

	if err := p.DB.Save(&product).Error; err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	p.setProductImage(ctx, product)

	return ctx.SendStatus(fiber.StatusOK)
}

//Delete - delete product
func (p *Product) Delete(ctx *fiber.Ctx) error {
	product, err := p.findProductByID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	p.DB.Unscoped().Delete(&product)
	p.removeImageProduct(ctx, product)
	return ctx.SendStatus(fiber.StatusNoContent)
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

	p.removeImageProduct(ctx, product)
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

func (p *Product) removeImageProduct(ctx *fiber.Ctx, product *models.Product) error {
	if product.Image != "" {
		product.Image = strings.Replace(product.Image, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + product.Image)
	}
	return nil
}
