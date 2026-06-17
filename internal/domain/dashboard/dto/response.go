package dto

type PopularBook struct {
	BookID    uint   `json:"book_id"`
	Title     string `json:"title"`
	LoanCount int64  `json:"loan_count"`
}

type LoanStats struct {
	Pending  int64 `json:"pending"`
	Borrowed int64 `json:"borrowed"`
	Returned int64 `json:"returned"`
	Rejected int64 `json:"rejected"`
	Overdue  int64 `json:"overdue"`
}

type Response struct {
	TotalUsers      int64         `json:"total_users"`
	TotalBooks      int64         `json:"total_books"`
	TotalCategories int64         `json:"total_categories"`
	TotalCopies     int64         `json:"total_copies"`
	AvailableCopies int64         `json:"available_copies"`
	BorrowedCopies  int64         `json:"borrowed_copies"`
	Loans           LoanStats     `json:"loans"`
	PopularBooks    []PopularBook `json:"popular_books"`
}

type UserStats struct {
	ActiveLoans   int64 `json:"active_loans"`
	TotalLoans    int64 `json:"total_loans"`
	ReturnedBooks int64 `json:"returned_books"`
	OverdueLoans  int64 `json:"overdue_loans"`
	PendingLoans  int64 `json:"pending_loans"`
}
