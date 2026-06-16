package book

import (
	"go-library/internal/domain/category"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title           string `json:"title" gorm:"type:varchar(200);not null"`
	Author          string `json:"author" gorm:"type:varchar(100);not null"`
	ISBN            string `json:"isbn" gorm:"type:varchar(20);uniqueIndex;not null"`
	Publisher       string `json:"publisher" gorm:"type:varchar(100)"`
	PublishedYear   int    `json:"published_year"`
	Description     string `json:"description" gorm:"type:text"`
	TotalCopies     int    `json:"total_copies" gorm:"not null;default:1"`
	AvailableCopies int    `json:"available_copies" gorm:"not null;default:1"`

	CategoryID uint              `json:"category_id" gorm:"not null"`
	Category   category.Category `json:"category" gorm:"foreignKey:CategoryID"`
}
