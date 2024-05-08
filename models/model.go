package models

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Notes struct {
	ID          string `json:"id" form:"ID"`
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	Content     string `json:"content" form:"content"`
	UserID      string `json:"user_id" form:"userID"`
}
