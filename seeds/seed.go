package seeds

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/bxcodec/faker/v3"
	"github.com/sing3demons/product-app/database"
	"github.com/sing3demons/product-app/models"
)

func Load() {
	db := database.GetDB()
	db.AutoMigrate(&models.Category{}, &models.Product{})

	var categories []models.Category
	err := db.Find(&categories).Error
	if len(categories) == 0 && err == nil {
		fmt.Println("Creating categories...")

		category := [...]string{"CPU", "GPU"}
		for i := 0; i < len(category); i++ {
			category := models.Category{
				Name: category[i],
			}

			categories = append(categories, category)
		}
		db.CreateInBatches(categories, len(category))
		fmt.Println("success")
	}

	numOfProducts := 100000
	products := make([]models.Product, numOfProducts)
	err = db.Find(&products).Limit(100).Error
	if len(products) == 0 && err == nil {
		fmt.Println("Creating products...")

		for i := 0; i < numOfProducts; i++ {
			product := models.Product{
				Name:       faker.Name(),
				Desc:       faker.Word(),
				Price:      rand.Intn(9999),
				Image:      "https://source.unsplash.com/random/300x200?" + strconv.Itoa(i),
				CategoryID: uint(rand.Intn(2) + 1),
			}
			products = append(products, product)
		}
		db.CreateInBatches(products, 1000)
		fmt.Println("success")

	}

}
