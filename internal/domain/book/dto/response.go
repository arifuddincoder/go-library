package dto

type CategoryInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	ID              uint         `json:"id"`
	Title           string       `json:"title"`
	Author          string       `json:"author"`
	ISBN            string       `json:"isbn"`
	Publisher       string       `json:"publisher,omitempty"`
	PublishedYear   int          `json:"published_year,omitempty"`
	Description     string       `json:"description,omitempty"`
	TotalCopies     int          `json:"total_copies"`
	AvailableCopies int          `json:"available_copies"`
	Category        CategoryInfo `json:"category"`
	CreatedAt       string       `json:"created_at,omitempty"`
}
