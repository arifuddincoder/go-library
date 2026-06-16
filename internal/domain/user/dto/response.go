package dto

type Response struct {
	ID        uint   `json:"id"`
	Name      string `json:"name" `
	Email     string `json:"email"  `
	Role      string `json:"role"`
	Token     string `json:"token,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}
