package user

import (
	"errors"
	"go-library/internal/query"

	"gorm.io/gorm"
)

type Repository interface {
	RegisterUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetAllUsers(p query.Params) ([]User, int64, error)
	DeleteUser(id uint) error
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

func (r repository) GetUserByEmail(email string) (*User, error) {
	var user User
	result := r.db.Where(&User{Email: email}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

func (r *repository) DeleteUser(id uint) error {
	result := r.db.Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *repository) GetAllUsers(p query.Params) ([]User, int64, error) {
	var users []User
	var total int64

	searchScope := query.Search(p.Search, "name", "email")
	allowedSort := map[string]bool{"name": true, "email": true, "created_at": true}

	r.db.Model(&User{}).
		Scopes(searchScope).
		Count(&total)

	result := r.db.
		Scopes(
			searchScope,
			query.Sort(p, allowedSort),
			query.Paginate(p),
		).
		Find(&users)

	if result.Error != nil {
		return nil, 0, result.Error
	}
	return users, total, nil
}
