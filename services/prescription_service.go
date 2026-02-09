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
		"INSERT INTO prescriptions (docId, userId, symptoms, link) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
		prescription.DocID, prescription.UserID, prescription.Symptoms, prescription.Link).Scan(&prescription.ID, &prescription.CreatedAt)
	return err
}

// GetPrescription retrieves a prescription by ID
func (s *PrescriptionService) GetPrescription(ctx context.Context, presID int64) (*models.Prescription, error) {
	prescription := &models.Prescription{}
	err := s.db.QueryRow(ctx,
		"SELECT id, created_at, docId, userId, symptoms, link FROM prescriptions WHERE id = $1",
		presID).Scan(&prescription.ID, &prescription.CreatedAt, &prescription.DocID, &prescription.UserID, &prescription.Symptoms, &prescription.Link)
	if err != nil {
		return nil, err
	}
	return prescription, nil
}

// GetUserPrescriptions retrieves all prescriptions for a user
func (s *PrescriptionService) GetUserPrescriptions(ctx context.Context, userID int64) ([]*models.Prescription, error) {
	rows, err := s.db.Query(ctx,
		"SELECT id, created_at, docId, userId, symptoms, link FROM prescriptions WHERE userId = $1 ORDER BY created_at DESC",
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prescriptions []*models.Prescription
	for rows.Next() {
		prescription := &models.Prescription{}
		if err := rows.Scan(&prescription.ID, &prescription.CreatedAt, &prescription.DocID, &prescription.UserID, &prescription.Symptoms, &prescription.Link); err != nil {
			return nil, err
		}
		prescriptions = append(prescriptions, prescription)
	}
	return prescriptions, rows.Err()
}
