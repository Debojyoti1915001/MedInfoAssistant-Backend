package models

import "time"

// Prescription represents a medical prescription
type Prescription struct {
	ID        int64     `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Symptoms  string    `db:"symptoms" json:"symptoms"`
	Link      string    `db:"link" json:"link"`
	UserID    int64     `db:"userId" json:"userId"`
	DocID     int64     `db:"docId" json:"docId"`
}
