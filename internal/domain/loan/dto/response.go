package dto

type BookInfo struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	ISBN  string `json:"isbn"`
}

type UserInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Response struct {
	ID         uint     `json:"id"`
	User       UserInfo `json:"user"`
	Book       BookInfo `json:"book"`
	BorrowedAt string   `json:"borrowed_at,omitempty"`
	DueDate    string   `json:"due_date,omitempty"`
	ReturnedAt string   `json:"returned_at,omitempty"`
	Status     string   `json:"status"`
	CreatedAt  string   `json:"created_at,omitempty"`
}
