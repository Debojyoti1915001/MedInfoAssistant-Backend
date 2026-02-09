package models

import "time"

// User represents a user in the database
type User struct {
	ID        int64     `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Name      string    `db:"name" json:"name"`
	PhnNumber string    `db:"phnNumber" json:"phnNumber"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password,omitempty"`
}

// LoginRequest represents user login credentials
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	PhnNumber string `json:"phnNumber"`
	Token     string `json:"token"`
}
