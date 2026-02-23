package entity

type Status string

const (
	StatusActive    Status = "active"
	StatusCompleted Status = "completed"
	StatusPending   Status = "pending"
)

type Task struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	UserID       int    `json:"user_id"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}
