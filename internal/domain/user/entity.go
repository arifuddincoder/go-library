package user

import (
	"go-library/internal/constants"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string         `json:"name" gorm:"type:varchar(100);not null"`
	Email    string         `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password string         `json:"_"  gorm:"type:varchar(100);not null"`
	Role     constants.Role `json:"role" gorm:"type:varchar(20);not null;default:user"`
}

func (u *User) hashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) checkPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
