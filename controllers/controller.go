package controllers

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type pagingResult struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	PrevPage  int   `json:"prevPage"`
	NextPage  int   `json:"nextPage"`
	Count     int64 `json:"count"`
	TotalPage int   `json:"totalPage"`
}

type pagination struct {
	ctx     *fiber.Ctx
	query   *gorm.DB
	records interface{}
}

func (p *pagination) paginate() *pagingResult {
	page, _ := strconv.Atoi(p.ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(p.ctx.Query("limit", "24"))

	ch := make(chan int64)
	go p.countRecords(ch)

	offset := (page - 1) * limit

	query := p.query.Preload("Category").Order("id desc")
	if category := p.ctx.Query("category"); category != "" {
		c, _ := strconv.Atoi(category)
		query = query.Where("category_id = ?", c)
	}

	query.Limit(limit).Offset(offset).Find(p.records)

	count := <-ch
	totalPage := int(math.Ceil(float64(count) / float64(limit)))

	var nextPage int
	if page == totalPage {
		nextPage = totalPage
	} else {
		nextPage = page + 1
	}

	return &pagingResult{
		Page:      page,
		Limit:     limit,
		PrevPage:  page - 1,
		NextPage:  nextPage,
		Count:     count,
		TotalPage: totalPage,
	}
}

func (p *pagination) countRecords(ch chan int64) {
	var count int64
	p.query.Model(p.records).Count(&count)
	ch <- count
}
