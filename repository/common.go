package repository

import (
	"device-management/model"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

func Paginate(page, perPage int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := page * perPage
		return db.Offset(offset).Limit(perPage)
	}
}

func Order(sort model.Sort) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(fmt.Sprintf("%s %s", sort.SortBy, strings.ToLower(sort.Direction)))
	}
}
