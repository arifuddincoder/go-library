package book

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	CreateBook(book *Book) error
	GetBookByID(id uint) (*Book, error)
	GetAllBooks() ([]Book, error)
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

func (r *repository) GetAllBooks() ([]Book, error) {
	var books []Book
	result := r.db.Preload("Category").Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}
	return books, nil
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
