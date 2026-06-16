package dto

type CreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}
