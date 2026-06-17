package user

import (
	"fmt"
	"go-library/internal/auth"
	"go-library/internal/constants"
	"go-library/internal/domain/user/dto"
	"go-library/internal/httpresponse"
	"go-library/internal/query"
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
		Role:  constants.RoleUser,
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

	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Name, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	response := dto.Response{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Role:         string(user.Role),
		Token:        token,
		RefreshToken: refreshToken,
		CreatedAt:    user.CreatedAt.String(),
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
		Role:  constants.RoleAdmin,
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

func (s *service) DeleteUser(id uint) error {
	return s.repo.DeleteUser(id)
}

func (s *service) GetAllUsers(p query.Params) (*httpresponse.Paginated, error) {
	users, total, err := s.repo.GetAllUsers(p)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.Response, 0, len(users))
	for _, user := range users {
		responses = append(responses, dto.Response{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt.String(),
		})
	}

	result := httpresponse.NewPaginated(responses, p.Page, p.Limit, total)
	return &result, nil
}

func (s *service) Refresh(refreshToken string) (*dto.Response, error) {
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.repo.GetUserByID(claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Name, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.Response{Token: token}, nil
}
