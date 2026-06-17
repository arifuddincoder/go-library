package loan

import (
	"time"

	"go-library/internal/domain/loan/dto"
	"go-library/internal/httpresponse"
	"go-library/internal/query"
)

const loanDuration = 14 * 24 * time.Hour

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func toResponse(l *Loan) *dto.Response {
	resp := &dto.Response{
		ID: l.ID,
		User: dto.UserInfo{
			ID:    l.User.ID,
			Name:  l.User.Name,
			Email: l.User.Email,
		},
		Book: dto.BookInfo{
			ID:    l.Book.ID,
			Title: l.Book.Title,
			ISBN:  l.Book.ISBN,
		},
		Status:    l.Status,
		CreatedAt: l.CreatedAt.String(),
	}
	if l.BorrowedAt != nil {
		resp.BorrowedAt = l.BorrowedAt.String()
	}
	if l.DueDate != nil {
		resp.DueDate = l.DueDate.String()
	}
	if l.ReturnedAt != nil {
		resp.ReturnedAt = l.ReturnedAt.String()
	}
	return resp
}

// user request kore — pending state e toiri hoy
func (s *service) RequestLoan(userID uint, req dto.RequestLoan) (*dto.Response, error) {
	// book exist kore kina check
	b, err := s.repo.GetBookByID(req.BookID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrBookNotFound
	}

	// copy ache kina check
	if b.AvailableCopies < 1 {
		return nil, ErrNoCopiesLeft
	}

	// already pending/borrowed ache kina check
	active, err := s.repo.HasActiveLoan(userID, req.BookID)
	if err != nil {
		return nil, err
	}
	if active {
		return nil, ErrAlreadyRequested
	}

	loan := Loan{
		UserID: userID,
		BookID: req.BookID,
		Status: StatusPending,
	}
	if err := s.repo.CreateLoan(&loan); err != nil {
		return nil, err
	}
	created, err := s.repo.GetLoanByID(loan.ID)
	if err != nil {
		return nil, err
	}
	return toResponse(created), nil
}

// admin accept kore — ekhane copy kome
func (s *service) ApproveLoan(loanID uint) (*dto.Response, error) {
	loan, err := s.repo.GetLoanByID(loanID)
	if err != nil {
		return nil, err
	}
	if loan == nil {
		return nil, ErrLoanNotFound
	}
	if loan.Status != StatusPending {
		return nil, ErrNotPending
	}

	now := time.Now()
	due := now.Add(loanDuration)
	loan.BorrowedAt = &now
	loan.DueDate = &due
	loan.Status = StatusBorrowed

	if err := s.repo.ApproveLoan(loan); err != nil {
		return nil, err
	}
	updated, err := s.repo.GetLoanByID(loan.ID)
	if err != nil {
		return nil, err
	}
	return toResponse(updated), nil
}

func (s *service) RejectLoan(loanID uint) (*dto.Response, error) {
	loan, err := s.repo.GetLoanByID(loanID)
	if err != nil {
		return nil, err
	}
	if loan == nil {
		return nil, ErrLoanNotFound
	}
	if loan.Status != StatusPending {
		return nil, ErrNotPending
	}

	loan.Status = StatusRejected
	if err := s.repo.UpdateLoan(loan); err != nil {
		return nil, err
	}
	return toResponse(loan), nil
}

// admin boi ferot ney — ekhane copy bare
func (s *service) ReturnLoan(loanID uint) (*dto.Response, error) {
	loan, err := s.repo.GetLoanByID(loanID)
	if err != nil {
		return nil, err
	}
	if loan == nil {
		return nil, ErrLoanNotFound
	}
	if loan.Status != StatusBorrowed {
		return nil, ErrNotBorrowed
	}

	now := time.Now()
	loan.ReturnedAt = &now
	loan.Status = StatusReturned

	if err := s.repo.ReturnLoan(loan); err != nil {
		return nil, err
	}
	return toResponse(loan), nil
}

func (s *service) paginate(loans []Loan, total int64, p query.Params) *httpresponse.Paginated {
	responses := make([]dto.Response, 0, len(loans))
	for i := range loans {
		responses = append(responses, *toResponse(&loans[i]))
	}
	result := httpresponse.NewPaginated(responses, p.Page, p.Limit, total)
	return &result
}

// present loan (borrowed)
func (s *service) GetMyActiveLoans(userID uint, p query.Params) (*httpresponse.Paginated, error) {
	loans, total, err := s.repo.GetLoansByUser(userID, []string{StatusBorrowed}, p)
	if err != nil {
		return nil, err
	}
	return s.paginate(loans, total, p), nil
}

// old loan (returned)
func (s *service) GetMyLoanHistory(userID uint, p query.Params) (*httpresponse.Paginated, error) {
	loans, total, err := s.repo.GetLoansByUser(userID, []string{StatusReturned}, p)
	if err != nil {
		return nil, err
	}
	return s.paginate(loans, total, p), nil
}

// admin — status diye filter kora jay (?status=pending)
func (s *service) GetAllLoans(statuses []string, p query.Params) (*httpresponse.Paginated, error) {
	loans, total, err := s.repo.GetAllLoans(statuses, p)
	if err != nil {
		return nil, err
	}
	return s.paginate(loans, total, p), nil
}
