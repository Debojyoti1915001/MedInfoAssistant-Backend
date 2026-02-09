package models

import "time"

// Doctor represents a medical doctor
type Doctor struct {
	ID         int64     `db:"id" json:"id"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	Accuracy   float64   `db:"accuracy" json:"accuracy"`
	Name       string    `db:"name" json:"name"`
	PhnNumber  string    `db:"phnNumber" json:"phnNumber"`
	Speciality string    `db:"speciality" json:"speciality"`
	Username   string    `db:"username" json:"username"`
	Email      string    `db:"email" json:"email"`
	Password   string    `db:"password" json:"password,omitempty"`
}

// DoctorCreateRequest represents doctor registration data (without accuracy)
type DoctorCreateRequest struct {
	Name       string `json:"name"`
	PhnNumber  string `json:"phnNumber"`
	Speciality string `json:"speciality"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

// DoctorLoginRequest represents doctor login credentials
type DoctorLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// DoctorLoginResponse represents the response after successful doctor login
type DoctorLoginResponse struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Username   string  `json:"username"`
	Speciality string  `json:"speciality"`
	Accuracy   float64 `json:"accuracy"`
	Token      string  `json:"token"`
}
