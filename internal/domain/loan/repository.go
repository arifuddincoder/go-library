package loan

import (
	"errors"
	"go-library/internal/domain/book"
	"go-library/internal/query"

	"gorm.io/gorm"
)

type Repository interface {
	CreateLoan(loan *Loan) error
	GetLoanByID(id uint) (*Loan, error)
	ApproveLoan(loan *Loan) error
	UpdateLoan(loan *Loan) error
	ReturnLoan(loan *Loan) error
	GetLoansByUser(userID uint, statuses []string, p query.Params) ([]Loan, int64, error)
	GetAllLoans(statuses []string, p query.Params) ([]Loan, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func statusScope(statuses []string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(statuses) == 0 {
			return db
		}
		return db.Where("status IN ?", statuses)
	}
}

func (r *repository) CreateLoan(loan *Loan) error {
	return r.db.Create(loan).Error
}

func (r *repository) GetLoanByID(id uint) (*Loan, error) {
	var loan Loan
	result := r.db.Preload("Book").Preload("User").First(&loan, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &loan, nil
}

// approve: boi available kina check kore, copy komay, loan save kore — sob ek transaction e
func (r *repository) ApproveLoan(loan *Loan) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var b book.Book
		if err := tx.First(&b, loan.BookID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBookNotFound
			}
			return err
		}
		if b.AvailableCopies < 1 {
			return ErrNoCopiesLeft
		}
		if err := tx.Model(&b).
			UpdateColumn("available_copies", gorm.Expr("available_copies - 1")).Error; err != nil {
			return err
		}
		return tx.Save(loan).Error
	})
}

// reject er jonno — sudhu loan save
func (r *repository) UpdateLoan(loan *Loan) error {
	return r.db.Save(loan).Error
}

// return: loan save kore copy barabe — ek transaction e
func (r *repository) ReturnLoan(loan *Loan) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(loan).Error; err != nil {
			return err
		}
		return tx.Model(&book.Book{}).
			Where("id = ?", loan.BookID).
			UpdateColumn("available_copies", gorm.Expr("available_copies + 1")).Error
	})
}

func (r *repository) GetLoansByUser(userID uint, statuses []string, p query.Params) ([]Loan, int64, error) {
	var loans []Loan
	var total int64

	r.db.Model(&Loan{}).
		Where("user_id = ?", userID).
		Scopes(statusScope(statuses)).
		Count(&total)

	allowedSort := map[string]bool{"created_at": true, "due_date": true}
	result := r.db.
		Preload("Book").Preload("User").
		Where("user_id = ?", userID).
		Scopes(statusScope(statuses), query.Sort(p, allowedSort), query.Paginate(p)).
		Find(&loans)

	if result.Error != nil {
		return nil, 0, result.Error
	}
	return loans, total, nil
}

func (r *repository) GetAllLoans(statuses []string, p query.Params) ([]Loan, int64, error) {
	var loans []Loan
	var total int64

	r.db.Model(&Loan{}).
		Scopes(statusScope(statuses)).
		Count(&total)

	allowedSort := map[string]bool{"created_at": true, "due_date": true}
	result := r.db.
		Preload("Book").Preload("User").
		Scopes(statusScope(statuses), query.Sort(p, allowedSort), query.Paginate(p)).
		Find(&loans)

	if result.Error != nil {
		return nil, 0, result.Error
	}
	return loans, total, nil
}
