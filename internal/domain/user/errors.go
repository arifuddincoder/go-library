package user

import "errors"

var (
	ErrorAlreadyExist     = errors.New("user with this email already exist")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid or expired refresh token")
)
