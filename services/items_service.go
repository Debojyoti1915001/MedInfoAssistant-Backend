package services

import (
	"context"
	"fmt"
	"strings"

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

// CreateItemsBulk creates multiple items for a prescription in a single query.
func (s *ItemsService) CreateItemsBulk(ctx context.Context, items []*models.Items) error {
	if len(items) == 0 {
		return nil
	}

	valueParts := make([]string, 0, len(items))
	args := make([]interface{}, 0, len(items)*5)
	for i, item := range items {
		base := i*5 + 1
		valueParts = append(valueParts, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", base, base+1, base+2, base+3, base+4))
		args = append(args, item.PresID, item.Name, item.Type, item.AIReasons, item.DocReason)
	}

	query := "INSERT INTO items (presId, name, type, aiReasons, docReason) VALUES " + strings.Join(valueParts, ", ")
	execArgs := make([]interface{}, 0, len(args)+1)
	execArgs = append(execArgs, pgx.QueryExecModeSimpleProtocol)
	execArgs = append(execArgs, args...)
	_, err := s.db.Exec(ctx, query, execArgs...)
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
