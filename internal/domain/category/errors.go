package category

import "errors"

var (
	ErrCategoryAlreadyExist = errors.New("category with this name already exist")
	ErrCategoryNotFound     = errors.New("category not found")
)
