package user

import "go-library/internal/domain/user/dto"

type service struct {
	repo Repository
	// jwtService auth.JWTService
}

func NewService(repo Repository) *service {
	return &service{repo}
}

func (s *service) RegisterUser(req dto.RegisterRequest) (*dto.Response, error) {
	user := User{
		Name:  req.Name,
		Email: req.Email,
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
		CreatedAt: user.CreatedAt.String(),
	}
	return &response, nil
}
