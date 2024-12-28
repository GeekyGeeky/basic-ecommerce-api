package models

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"` // Don't expose the password
	IsAdmin  bool   `json:"is_admin"`
}
