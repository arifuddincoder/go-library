package user

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	RegisterUser(user *User) error
	// GetUserByEmail(email string) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) RegisterUser(user *User) error {
	result := r.db.Create(user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrorAlreadyExist
		}

		return result.Error
	}

	return nil
}
