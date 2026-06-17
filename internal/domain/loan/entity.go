package loan

import (
	"time"

	"go-library/internal/domain/book"
	"go-library/internal/domain/user"

	"gorm.io/gorm"
)

const (
	StatusPending  = "pending"
	StatusBorrowed = "borrowed"
	StatusReturned = "returned"
	StatusRejected = "rejected"
)

type Loan struct {
	gorm.Model
	UserID     uint       `json:"user_id" gorm:"not null"`
	BookID     uint       `json:"book_id" gorm:"not null"`
	BorrowedAt *time.Time `json:"borrowed_at"`
	DueDate    *time.Time `json:"due_date"`
	ReturnedAt *time.Time `json:"returned_at"`
	Status     string     `json:"status" gorm:"type:varchar(20);not null;default:pending"`

	User user.User `json:"user" gorm:"foreignKey:UserID"`
	Book book.Book `json:"book" gorm:"foreignKey:BookID"`
}
