package category

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	CreateCategory(category *Category) error
	GetCategoryByName(name string) (*Category, error)
	GetCategoryByID(id uint) (*Category, error)
	GetAllCategories() ([]Category, error)
	DeleteCategory(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateCategory(category *Category) error {
	result := r.db.Create(category)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrCategoryAlreadyExist
		}
		return result.Error
	}
	return nil
}

func (r *repository) GetCategoryByName(name string) (*Category, error) {
	var category Category
	result := r.db.Where(&Category{Name: name}).First(&category)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &category, nil
}

func (r *repository) GetCategoryByID(id uint) (*Category, error) {
	var category Category
	result := r.db.First(&category, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &category, nil
}

func (r *repository) GetAllCategories() ([]Category, error) {
	var categories []Category
	result := r.db.Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}
	return categories, nil
}

func (r *repository) DeleteCategory(id uint) error {
	result := r.db.Delete(&Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCategoryNotFound
	}
	return nil
}
