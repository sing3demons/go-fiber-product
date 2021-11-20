package models

import "gorm.io/gorm"

//Product - model
type Product struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Desc  string `gorm:"not null"`
	Price int    `gorm:"not null"`
	Image string `gorm:"not null"`
	CategoryID uint
	Category   Category
}
