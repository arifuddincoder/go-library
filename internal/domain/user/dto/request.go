package dto

import "gorm.io/gorm"

type RegisterRequest struct {
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(100);not null"`
	Email    string `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password string `json:"password" gorm:"type:varchar(100);not null"`
	Role     string `json:"role" gorm:"type:varchar(20);not null;default:user"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
