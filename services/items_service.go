package services

import (
	"context"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/models"
	"github.com/jackc/pgx/v5"
)

type ItemsService struct {
	db *pgx.Conn
}

func NewItemsService(db *pgx.Conn) *ItemsService {
	return &ItemsService{db: db}
}

// CreateItem creates a new item in a prescription
func (s *ItemsService) CreateItem(ctx context.Context, item *models.Items) error {
	err := s.db.QueryRow(ctx,
		"INSERT INTO items (presId, name, type, aiReasons, docReason) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at",
		item.PresID, item.Name, item.Type, item.AIReasons, item.DocReason).Scan(&item.ID, &item.CreatedAt)
	return err
}

// GetItem retrieves an item by ID
func (s *ItemsService) GetItem(ctx context.Context, itemID int64) (*models.Items, error) {
	item := &models.Items{}
	err := s.db.QueryRow(ctx,
		"SELECT id, created_at, name, type, aiReasons, docReason, presId FROM items WHERE id = $1",
		itemID).Scan(&item.ID, &item.CreatedAt, &item.Name, &item.Type, &item.AIReasons, &item.DocReason, &item.PresID)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// GetPrescriptionItems retrieves all items for a prescription
func (s *ItemsService) GetPrescriptionItems(ctx context.Context, presID int64) ([]*models.Items, error) {
	rows, err := s.db.Query(ctx,
		"SELECT id, created_at, name, type, aiReasons, docReason, presId FROM items WHERE presId = $1",
		presID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.Items
	for rows.Next() {
		item := &models.Items{}
		if err := rows.Scan(&item.ID, &item.CreatedAt, &item.Name, &item.Type, &item.AIReasons, &item.DocReason, &item.PresID); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
