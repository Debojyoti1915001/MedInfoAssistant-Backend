package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/models"
	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/services"
	"github.com/jackc/pgx/v5"
)

// CreatePrescriptionHandler creates a new prescription
func CreatePrescriptionHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var prescription models.Prescription
		if err := json.NewDecoder(r.Body).Decode(&prescription); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		presService := services.NewPrescriptionService(db)
		if err := presService.CreatePrescription(context.Background(), &prescription); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(prescription)
	}
}

// GetPrescriptionHandler returns a specific prescription
func GetPrescriptionHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Prescription ID is required", http.StatusBadRequest)
			return
		}

		presID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Prescription ID", http.StatusBadRequest)
			return
		}

		presService := services.NewPrescriptionService(db)
		prescription, err := presService.GetPrescription(context.Background(), presID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prescription)
	}
}

// GetUserPrescriptionsHandler returns all prescriptions for a user
func GetUserPrescriptionsHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userIDStr := r.URL.Query().Get("userId")
		if userIDStr == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid User ID", http.StatusBadRequest)
			return
		}

		presService := services.NewPrescriptionService(db)
		prescriptions, err := presService.GetUserPrescriptions(context.Background(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prescriptions)
	}
}
