package dashboard

import (
	"time"

	"go-library/internal/domain/book"
	"go-library/internal/domain/category"
	"go-library/internal/domain/dashboard/dto"
	"go-library/internal/domain/loan"
	"go-library/internal/domain/user"

	"gorm.io/gorm"
)

type Repository interface {
	GetStats() (*dto.Response, error)
	GetUserStats(userID uint) (*dto.UserStats, error)
}
type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetStats() (*dto.Response, error) {
	var res dto.Response

	r.db.Model(&user.User{}).Count(&res.TotalUsers)
	r.db.Model(&category.Category{}).Count(&res.TotalCategories)
	r.db.Model(&book.Book{}).Count(&res.TotalBooks)

	r.db.Model(&book.Book{}).Select("COALESCE(SUM(total_copies),0)").Scan(&res.TotalCopies)
	r.db.Model(&book.Book{}).Select("COALESCE(SUM(available_copies),0)").Scan(&res.AvailableCopies)
	res.BorrowedCopies = res.TotalCopies - res.AvailableCopies

	r.db.Model(&loan.Loan{}).Where("status = ?", loan.StatusPending).Count(&res.Loans.Pending)
	r.db.Model(&loan.Loan{}).Where("status = ?", loan.StatusBorrowed).Count(&res.Loans.Borrowed)
	r.db.Model(&loan.Loan{}).Where("status = ?", loan.StatusReturned).Count(&res.Loans.Returned)
	r.db.Model(&loan.Loan{}).Where("status = ?", loan.StatusRejected).Count(&res.Loans.Rejected)
	r.db.Model(&loan.Loan{}).
		Where("status = ? AND due_date < ?", loan.StatusBorrowed, time.Now()).
		Count(&res.Loans.Overdue)

	r.db.Model(&loan.Loan{}).
		Select("loans.book_id, books.title as title, COUNT(loans.id) as loan_count").
		Joins("JOIN books ON books.id = loans.book_id").
		Where("loans.status IN ?", []string{loan.StatusBorrowed, loan.StatusReturned}).
		Group("loans.book_id, books.title").
		Order("loan_count DESC").
		Limit(5).
		Scan(&res.PopularBooks)

	return &res, nil
}

func (r *repository) GetUserStats(userID uint) (*dto.UserStats, error) {
	var res dto.UserStats
	base := r.db.Model(&loan.Loan{}).Where("user_id = ?", userID)

	base.Session(&gorm.Session{}).Count(&res.TotalLoans)
	base.Session(&gorm.Session{}).Where("status = ?", loan.StatusBorrowed).Count(&res.ActiveLoans)
	base.Session(&gorm.Session{}).Where("status = ?", loan.StatusReturned).Count(&res.ReturnedBooks)
	base.Session(&gorm.Session{}).Where("status = ?", loan.StatusPending).Count(&res.PendingLoans)
	base.Session(&gorm.Session{}).
		Where("status = ? AND due_date < ?", loan.StatusBorrowed, time.Now()).
		Count(&res.OverdueLoans)

	return &res, nil
}
