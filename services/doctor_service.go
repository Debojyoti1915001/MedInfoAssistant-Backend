package services

import (
	"context"
	"errors"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/models"
	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/utils"
	"github.com/jackc/pgx/v5"
)

type DoctorService struct {
	db *pgx.Conn
}

func NewDoctorService(db *pgx.Conn) *DoctorService {
	return &DoctorService{db: db}
}

// CreateDoctor creates a new doctor with hashed password
func (s *DoctorService) CreateDoctor(ctx context.Context, doctor *models.Doctor) error {
	// Hash the password
	hashedPassword, err := utils.HashPassword(doctor.Password)
	if err != nil {
		return err
	}

	err = s.db.QueryRow(ctx,
		"INSERT INTO doctors (accuracy, name, phnNumber, speciality, username, email, password) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at",
		doctor.Accuracy, doctor.Name, doctor.PhnNumber, doctor.Speciality, doctor.Username, doctor.Email, hashedPassword).Scan(&doctor.ID, &doctor.CreatedAt)
	return err
}

// CreateDoctorWithRequest creates a new doctor from DoctorCreateRequest (accuracy defaults to 0.0)
func (s *DoctorService) CreateDoctorWithRequest(ctx context.Context, req *models.DoctorCreateRequest) (*models.Doctor, error) {
	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	doctor := &models.Doctor{
		Accuracy:   0.0, // Initialize accuracy to 0.0 for new doctors
		Name:       req.Name,
		PhnNumber:  req.PhnNumber,
		Speciality: req.Speciality,
		Username:   req.Username,
		Email:      req.Email,
	}

	err = s.db.QueryRow(ctx,
		"INSERT INTO doctors (accuracy, name, phnNumber, speciality, username, email, password) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at",
		doctor.Accuracy, doctor.Name, doctor.PhnNumber, doctor.Speciality, doctor.Username, doctor.Email, hashedPassword).Scan(&doctor.ID, &doctor.CreatedAt)
	if err != nil {
		return nil, err
	}

	return doctor, nil
}

// GetDoctor retrieves a doctor by ID
func (s *DoctorService) GetDoctor(ctx context.Context, docID int64) (*models.Doctor, error) {
	doctor := &models.Doctor{}
	err := s.db.QueryRow(ctx,
		"SELECT id, created_at, accuracy, name, phnNumber, speciality, username, email, password FROM doctors WHERE id = $1",
		docID).Scan(&doctor.ID, &doctor.CreatedAt, &doctor.Accuracy, &doctor.Name, &doctor.PhnNumber, &doctor.Speciality, &doctor.Username, &doctor.Email, &doctor.Password)
	if err != nil {
		return nil, err
	}
	return doctor, nil
}

// GetAllDoctors retrieves all doctors
func (s *DoctorService) GetAllDoctors(ctx context.Context) ([]*models.Doctor, error) {
	rows, err := s.db.Query(ctx,
		"SELECT id, created_at, accuracy, name, phnNumber, speciality, username, email, password FROM doctors ORDER BY accuracy DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []*models.Doctor
	for rows.Next() {
		doctor := &models.Doctor{}
		if err := rows.Scan(&doctor.ID, &doctor.CreatedAt, &doctor.Accuracy, &doctor.Name, &doctor.PhnNumber, &doctor.Speciality, &doctor.Username, &doctor.Email, &doctor.Password); err != nil {
			return nil, err
		}
		doctors = append(doctors, doctor)
	}
	return doctors, rows.Err()
}

// GetDoctorByUsername retrieves a doctor by username
func (s *DoctorService) GetDoctorByUsername(ctx context.Context, username string) (*models.Doctor, error) {
	doctor := &models.Doctor{}
	err := s.db.QueryRow(ctx,
		"SELECT id, created_at, accuracy, name, phnNumber, speciality, username, email, password FROM doctors WHERE username = $1",
		username).Scan(&doctor.ID, &doctor.CreatedAt, &doctor.Accuracy, &doctor.Name, &doctor.PhnNumber, &doctor.Speciality, &doctor.Username, &doctor.Email, &doctor.Password)
	if err != nil {
		return nil, err
	}
	return doctor, nil
}

// LoginDoctor authenticates a doctor by email and password
func (s *DoctorService) LoginDoctor(ctx context.Context, email, password string) (*models.Doctor, error) {
	doctor := &models.Doctor{}
	err := s.db.QueryRow(ctx, "SELECT id, created_at, accuracy, name, phnNumber, speciality, username, email, password FROM doctors WHERE email = $1", email).Scan(&doctor.ID, &doctor.CreatedAt, &doctor.Accuracy, &doctor.Name, &doctor.PhnNumber, &doctor.Speciality, &doctor.Username, &doctor.Email, &doctor.Password)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify the password
	if !utils.VerifyPassword(doctor.Password, password) {
		return nil, errors.New("invalid email or password")
	}

	return doctor, nil
}
