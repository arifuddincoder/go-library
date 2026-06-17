package category

import (
	"go-library/internal/domain/category/dto"
	"go-library/internal/httpresponse"
	"go-library/internal/query"
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

func (s *service) GetAllCategories(p query.Params) (*httpresponse.Paginated, error) {
	categories, total, err := s.repo.GetAllCategories(p)
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

	result := httpresponse.NewPaginated(responses, p.Page, p.Limit, total)
	return &result, nil
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
