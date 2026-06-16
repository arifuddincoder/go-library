package dto

type RegisterRequest struct {
	Title         string `json:"title" validate:"required"`
	Author        string `json:"author" validate:"required"`
	ISBN          string `json:"isbn" validate:"required"`
	Publisher     string `json:"publisher"`
	PublishedYear int    `json:"published_year"`
	CategoryID    uint   `json:"category_id" validate:"required"`
	Description   string `json:"description"`
	TotalCopies   int    `json:"total_copies" validate:"required,min=1"`
}
