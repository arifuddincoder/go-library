package category

import (
	"errors"
	"go-library/internal/query"

	"gorm.io/gorm"
)

type Repository interface {
	CreateCategory(category *Category) error
	GetCategoryByName(name string) (*Category, error)
	GetCategoryByID(id uint) (*Category, error)
	GetAllCategories(p query.Params) ([]Category, int64, error)
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

func (r *repository) GetAllCategories(p query.Params) ([]Category, int64, error) {
	var categories []Category
	var total int64

	searchScope := query.Search(p.Search, "name", "description")
	allowedSort := map[string]bool{"name": true, "created_at": true}

	r.db.Model(&Category{}).
		Scopes(searchScope).
		Count(&total)

	result := r.db.
		Scopes(
			searchScope,
			query.Sort(p, allowedSort),
			query.Paginate(p),
		).
		Find(&categories)

	if result.Error != nil {
		return nil, 0, result.Error
	}
	return categories, total, nil
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
