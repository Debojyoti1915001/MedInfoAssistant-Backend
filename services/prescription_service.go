package services

import (
	"context"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/models"
	"github.com/jackc/pgx/v5"
)

type PrescriptionService struct {
	db *pgx.Conn
}

func NewPrescriptionService(db *pgx.Conn) *PrescriptionService {
	return &PrescriptionService{db: db}
}

// CreatePrescription creates a new prescription
func (s *PrescriptionService) CreatePrescription(ctx context.Context, prescription *models.Prescription) error {
	err := s.db.QueryRow(ctx,
		`INSERT INTO prescriptions (docId, userId, symptoms, link, seenByPatient)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at, COALESCE(seenByPatient, FALSE)`,
		prescription.DocID, prescription.UserID, prescription.Symptoms, prescription.Link, false).
		Scan(&prescription.ID, &prescription.CreatedAt, &prescription.SeenByPatient)
	return err
}

// UpdatePrescriptionLink updates only the storage link for a prescription.
func (s *PrescriptionService) UpdatePrescriptionLink(ctx context.Context, presID int64, link string) error {
	_, err := s.db.Exec(ctx, "UPDATE prescriptions SET link = $2 WHERE id = $1", presID, link)
	return err
}

// GetPrescription retrieves a prescription by ID
func (s *PrescriptionService) GetPrescription(ctx context.Context, presID int64) (*models.Prescription, error) {
	prescription := &models.Prescription{}
	err := s.db.QueryRow(ctx,
		"SELECT id, created_at, docId, userId, symptoms, link, COALESCE(seenByPatient, FALSE) FROM prescriptions WHERE id = $1",
		presID).Scan(
		&prescription.ID,
		&prescription.CreatedAt,
		&prescription.DocID,
		&prescription.UserID,
		&prescription.Symptoms,
		&prescription.Link,
		&prescription.SeenByPatient,
	)
	if err != nil {
		return nil, err
	}
	return prescription, nil
}

// GetUserPrescriptions retrieves all prescriptions for a user
func (s *PrescriptionService) GetUserPrescriptions(ctx context.Context, userID int64) ([]*models.Prescription, error) {
	rows, err := s.db.Query(ctx,
		"SELECT id, created_at, docId, userId, symptoms, link, COALESCE(seenByPatient, FALSE) FROM prescriptions WHERE userId = $1 ORDER BY created_at DESC",
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prescriptions []*models.Prescription
	for rows.Next() {
		prescription := &models.Prescription{}
		if err := rows.Scan(
			&prescription.ID,
			&prescription.CreatedAt,
			&prescription.DocID,
			&prescription.UserID,
			&prescription.Symptoms,
			&prescription.Link,
			&prescription.SeenByPatient,
		); err != nil {
			return nil, err
		}
		prescriptions = append(prescriptions, prescription)
	}
	return prescriptions, rows.Err()
}

// GetDoctorPrescriptions retrieves all prescriptions for a doctor
func (s *PrescriptionService) GetDoctorPrescriptions(ctx context.Context, docID int64) ([]*models.Prescription, error) {
	rows, err := s.db.Query(ctx,
		"SELECT id, created_at, docId, userId, symptoms, link, COALESCE(seenByPatient, FALSE) FROM prescriptions WHERE docId = $1 ORDER BY created_at DESC",
		docID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prescriptions []*models.Prescription
	for rows.Next() {
		prescription := &models.Prescription{}
		if err := rows.Scan(
			&prescription.ID,
			&prescription.CreatedAt,
			&prescription.DocID,
			&prescription.UserID,
			&prescription.Symptoms,
			&prescription.Link,
			&prescription.SeenByPatient,
		); err != nil {
			return nil, err
		}
		prescriptions = append(prescriptions, prescription)
	}
	return prescriptions, rows.Err()
}

// UpdatePrescriptionSeenByPatient updates the seenByPatient status for a prescription by ID.
func (s *PrescriptionService) UpdatePrescriptionSeenByPatient(ctx context.Context, presID int64, seenByPatient bool) (*models.Prescription, error) {
	prescription := &models.Prescription{}
	err := s.db.QueryRow(ctx,
		`UPDATE prescriptions
		 SET seenByPatient = $2
		 WHERE id = $1
		 RETURNING id, created_at, docId, userId, symptoms, link, COALESCE(seenByPatient, FALSE)`,
		presID, seenByPatient,
	).Scan(
		&prescription.ID,
		&prescription.CreatedAt,
		&prescription.DocID,
		&prescription.UserID,
		&prescription.Symptoms,
		&prescription.Link,
		&prescription.SeenByPatient,
	)
	if err != nil {
		return nil, err
	}
	return prescription, nil
}
