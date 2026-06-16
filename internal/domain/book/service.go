package book

import "go-library/internal/domain/book/dto"

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

func (s *service) GetAllBooks() ([]dto.Response, error) {
	books, err := s.repo.GetAllBooks()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.Response, 0, len(books))
	for i := range books {
		responses = append(responses, *toResponse(&books[i]))
	}
	return responses, nil
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
