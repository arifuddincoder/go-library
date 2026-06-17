package book

import (
	"errors"
	"go-library/internal/query"

	"gorm.io/gorm"
)

type Repository interface {
	CreateBook(book *Book) error
	GetBookByID(id uint) (*Book, error)
	GetAllBooks(p query.Params) ([]Book, int64, error)
	UpdateBook(book *Book) error
	DeleteBook(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateBook(book *Book) error {
	result := r.db.Create(book)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrBookAlreadyExist
		}
		return result.Error
	}
	return nil
}

func (r *repository) GetBookByID(id uint) (*Book, error) {
	var book Book
	result := r.db.Preload("Category").First(&book, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &book, nil
}

func (r *repository) GetAllBooks(p query.Params) ([]Book, int64, error) {
	var books []Book
	var total int64

	// search column gulo specify koro
	searchScope := query.Search(p.Search, "title", "author", "isbn")
	allowedSort := map[string]bool{"title": true, "author": true, "created_at": true}

	// total count (pagination chara, but search soho)
	r.db.Model(&Book{}).
		Scopes(searchScope).
		Count(&total)

	// actual data (sob scope chain kore)
	result := r.db.
		Preload("Category").
		Scopes(
			searchScope,
			query.Sort(p, allowedSort),
			query.Paginate(p),
		).
		Find(&books)

	if result.Error != nil {
		return nil, 0, result.Error
	}
	return books, total, nil
}

func (r *repository) UpdateBook(book *Book) error {
	result := r.db.Save(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *repository) DeleteBook(id uint) error {
	result := r.db.Delete(&Book{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
