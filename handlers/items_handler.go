package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/models"
	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/services"
	"github.com/jackc/pgx/v5"
)

type updateItemDocReasonRequest struct {
	DocReason string `json:"docReason"`
}

// CreateItemHandler creates a new item in a prescription
func CreateItemHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var item models.Items
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		itemService := services.NewItemsService(db)
		if err := itemService.CreateItem(context.Background(), &item); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(item)
	}
}

// GetItemHandler returns a specific item
func GetItemHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Item ID is required", http.StatusBadRequest)
			return
		}

		itemID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Item ID", http.StatusBadRequest)
			return
		}

		itemService := services.NewItemsService(db)
		item, err := itemService.GetItem(context.Background(), itemID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item)
	}
}

// GetPrescriptionItemsHandler returns all items for a prescription
func GetPrescriptionItemsHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		presIDStr := r.URL.Query().Get("presId")
		if presIDStr == "" {
			http.Error(w, "Prescription ID is required", http.StatusBadRequest)
			return
		}

		presID, err := strconv.ParseInt(presIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Prescription ID", http.StatusBadRequest)
			return
		}

		itemService := services.NewItemsService(db)
		items, err := itemService.GetPrescriptionItems(context.Background(), presID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

// UpdateItemDocReasonHandler updates docReason for a specific item.
// Query param: id
// Body: {"docReason":"..."}
func UpdateItemDocReasonHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Item ID is required", http.StatusBadRequest)
			return
		}

		itemID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Item ID", http.StatusBadRequest)
			return
		}

		var req updateItemDocReasonRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		docReason := strings.TrimSpace(req.DocReason)

		itemService := services.NewItemsService(db)
		updatedItem, err := itemService.UpdateItemDocReason(context.Background(), itemID, docReason)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "item not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedItem)
	}
}
