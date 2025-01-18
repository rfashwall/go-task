package models

// User represents a user in our system.
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"` // Note: In a real application, never store passwords in plain text!
}
