package loan

import "errors"

var (
	ErrLoanNotFound = errors.New("loan not found")
	ErrBookNotFound = errors.New("book not found")
	ErrNoCopiesLeft = errors.New("no available copies left for this book")
	ErrNotPending   = errors.New("loan is not in pending state")
	ErrNotBorrowed  = errors.New("loan is not in borrowed state")
)
