package category

import (
	"go-library/internal/domain/category/dto"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

func (s *service) CreateCategory(req dto.CreateRequest) (*dto.Response, error) {
	existing, err := s.repo.GetCategoryByName(req.Name)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, ErrCategoryAlreadyExist
	}

	category := Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.CreateCategory(&category); err != nil {
		return nil, err
	}

	response := dto.Response{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt.String(),
	}
	return &response, nil
}

func (s *service) GetAllCategories() ([]dto.Response, error) {
	categories, err := s.repo.GetAllCategories()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.Response, 0, len(categories))
	for _, category := range categories {
		responses = append(responses, dto.Response{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			CreatedAt:   category.CreatedAt.String(),
		})
	}
	return responses, nil
}

func (s *service) DeleteCategory(id uint) error {
	existing, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrCategoryNotFound
	}

	return s.repo.DeleteCategory(id)
}
