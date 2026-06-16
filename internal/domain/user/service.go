package user

import (
	"fmt"
	"go-library/internal/auth"
	"go-library/internal/domain/user/dto"
)

type service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) *service {
	return &service{repo, jwtService}
}

func (s *service) RegisterUser(req dto.RegisterRequest) (*dto.Response, error) {

	existing, err := s.repo.GetUserByEmail(req.Email)

	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, ErrorAlreadyExist
	}

	user := User{
		Name:  req.Name,
		Email: req.Email,
		Role:  RoleUser,
	}

	if err := user.hashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.repo.RegisterUser(&user); err != nil {
		return nil, err
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt.String(),
	}
	return &response, nil
}

func (s *service) LoginUser(req dto.LoginRequest) (*dto.Response, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// check password
	err = user.checkPassword(req.Password)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// generate token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Name, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		Token:     token,
		CreatedAt: user.CreatedAt.String(),
	}

	return &response, nil
}

func (s *service) CreateAdmin(req dto.RegisterRequest) (*dto.Response, error) {
	existing, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, ErrorAlreadyExist
	}

	user := User{
		Name:  req.Name,
		Email: req.Email,
		Role:  RoleAdmin,
	}

	if err := user.hashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.repo.RegisterUser(&user); err != nil {
		return nil, err
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt.String(),
	}
	return &response, nil
}
