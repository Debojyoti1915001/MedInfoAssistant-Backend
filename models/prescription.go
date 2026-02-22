package models

import "time"

// Prescription represents a medical prescription
type Prescription struct {
	ID            int64     `db:"id" json:"id"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	Symptoms      string    `db:"symptoms" json:"symptoms"`
	Link          string    `db:"link" json:"link"`
	UserID        int64     `db:"userId" json:"userId"`
	DocID         int64     `db:"docId" json:"docId"`
	SeenByPatient bool      `db:"seenByPatient" json:"seenByPatient"`
}
type Test struct {
	Name       string  `json:"name"`
	Reason1    string  `json:"reason1"`
	Precision1 float64 `json:"precision1"`
	Reason2    string  `json:"reason2"`
	Precision2 float64 `json:"precision2"`
	Reason3    string  `json:"reason3"`
	Precision3 float64 `json:"precision3"`
}

type Medicine struct {
	Name         string  `json:"name"`
	Description1 string  `json:"description1"`
	Precision1   float64 `json:"precision1"`
	Description2 string  `json:"description2"`
	Precision2   float64 `json:"precision2"`
	Description3 string  `json:"description3"`
	Precision3   float64 `json:"precision3"`
	Price        float64 `json:"price"`
}

type AIResponse struct {
	Tests     map[string]Test     `json:"tests"`
	Medicines map[string]Medicine `json:"medicines"`
}

type AIRequest struct {
	File             string `json:"file"`
	Symptoms         string `json:"symptoms"`
	DoctorSpeciality string `json:"doctor_speciality"`
}
