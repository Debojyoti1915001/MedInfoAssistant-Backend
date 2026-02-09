package models

import "time"

// Items represents medical items (medicines or tests) in a prescription
type Items struct {
	ID        int64     `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Name      string    `db:"name" json:"name"`
	Type      string    `db:"type" json:"type"`
	AIReasons string    `db:"aiReasons" json:"aiReasons"`
	DocReason string    `db:"docReason" json:"docReason"`
	PresID    int64     `db:"presId" json:"presId"`
}
