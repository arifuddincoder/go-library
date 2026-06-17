package dto

type RequestLoan struct {
	BookID uint `json:"book_id" validate:"required"`
}
