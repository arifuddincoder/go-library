package book

import "errors"

var (
	ErrBookAlreadyExist   = errors.New("book with this ISBN already exist")
	ErrBookNotFound       = errors.New("book not found")
	ErrCategoryNotFound   = errors.New("category not found for the given id")
	ErrNoCopiesLeft       = errors.New("no available copies left for this book")
	ErrInvalidTotalCopies = errors.New("total copies cannot be less than borrowed copies")
)
