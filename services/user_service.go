package services

import (
	"context"
	"errors"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/models"
	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/utils"
	"github.com/jackc/pgx/v5"
)

type UserService struct {
	db *pgx.Conn
}

func NewUserService(db *pgx.Conn) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a new user in the database with hashed password
func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	err = s.db.QueryRow(ctx,
		"INSERT INTO users (name, phnNumber, email, password) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
		user.Name, user.PhnNumber, user.Email, hashedPassword).Scan(&user.ID, &user.CreatedAt)
	return err
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(ctx, "SELECT id, created_at, name, phnNumber, email, password FROM users WHERE id = $1", id).Scan(&user.ID, &user.CreatedAt, &user.Name, &user.PhnNumber, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetAllUsers retrieves all users from the database
func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	rows, err := s.db.Query(ctx, "SELECT id, created_at, name, phnNumber, email, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.CreatedAt, &user.Name, &user.PhnNumber, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// LoginUser authenticates a user by email and password
func (s *UserService) LoginUser(ctx context.Context, email, password string) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(ctx, "SELECT id, created_at, name, phnNumber, email, password FROM users WHERE email = $1", email).Scan(&user.ID, &user.CreatedAt, &user.Name, &user.PhnNumber, &user.Email, &user.Password)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify the password
	if !utils.VerifyPassword(user.Password, password) {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
