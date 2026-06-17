package book

import (
	"go-library/internal/domain/book/dto"
	"go-library/internal/httpresponse"
	"go-library/internal/query"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func toResponse(b *Book) *dto.Response {
	return &dto.Response{
		ID:              b.ID,
		Title:           b.Title,
		Author:          b.Author,
		ISBN:            b.ISBN,
		Publisher:       b.Publisher,
		PublishedYear:   b.PublishedYear,
		Description:     b.Description,
		TotalCopies:     b.TotalCopies,
		AvailableCopies: b.AvailableCopies,
		Category: dto.CategoryInfo{
			ID:   b.Category.ID,
			Name: b.Category.Name,
		},
		CreatedAt: b.CreatedAt.String(),
	}
}

func (s *service) CreateBook(req dto.RegisterRequest) (*dto.Response, error) {
	book := Book{
		Title:           req.Title,
		Author:          req.Author,
		ISBN:            req.ISBN,
		Publisher:       req.Publisher,
		PublishedYear:   req.PublishedYear,
		Description:     req.Description,
		TotalCopies:     req.TotalCopies,
		AvailableCopies: req.TotalCopies,
		CategoryID:      req.CategoryID,
	}

	if err := s.repo.CreateBook(&book); err != nil {
		return nil, err
	}

	created, err := s.repo.GetBookByID(book.ID)
	if err != nil {
		return nil, err
	}
	return toResponse(created), nil
}

func (s *service) GetBookByID(id uint) (*dto.Response, error) {
	book, err := s.repo.GetBookByID(id)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, ErrBookNotFound
	}
	return toResponse(book), nil
}

func (s *service) GetAllBooks(p query.Params) (*httpresponse.Paginated, error) {
	books, total, err := s.repo.GetAllBooks(p)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.Response, 0, len(books))
	for i := range books {
		responses = append(responses, *toResponse(&books[i]))
	}

	result := httpresponse.NewPaginated(responses, p.Page, p.Limit, total)
	return &result, nil
}

func (s *service) UpdateBook(id uint, req dto.UpdateRequest) (*dto.Response, error) {
	existing, err := s.repo.GetBookByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrBookNotFound
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Author != nil {
		existing.Author = *req.Author
	}
	if req.Publisher != nil {
		existing.Publisher = *req.Publisher
	}
	if req.PublishedYear != nil {
		existing.PublishedYear = *req.PublishedYear
	}
	if req.CategoryID != nil {
		existing.CategoryID = *req.CategoryID
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.TotalCopies != nil {
		// koto ta borrowed ache seita thik rekhe available adjust kori
		borrowed := existing.TotalCopies - existing.AvailableCopies
		newAvailable := *req.TotalCopies - borrowed
		if newAvailable < 0 {
			return nil, ErrInvalidTotalCopies
		}
		existing.TotalCopies = *req.TotalCopies
		existing.AvailableCopies = newAvailable
	}

	if err := s.repo.UpdateBook(existing); err != nil {
		return nil, err
	}

	// category name soho fresh data fetch kori
	updated, err := s.repo.GetBookByID(id)
	if err != nil {
		return nil, err
	}
	return toResponse(updated), nil
}

func (s *service) DeleteBook(id uint) error {
	book, err := s.repo.GetBookByID(id)
	if err != nil {
		return err
	}
	if book == nil {
		return ErrBookNotFound
	}
	return s.repo.DeleteBook(id)
}
