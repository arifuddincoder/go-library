package category

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string `json:"name" gorm:"type:varchar(50);uniqueIndex;not null"`
	Description string `json:"description" gorm:"type:varchar(255)"`
}
