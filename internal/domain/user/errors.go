package user

import "errors"

var (
	ErrorAlreadyExist     = errors.New("user with this email already exist")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
