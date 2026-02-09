package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/models"
	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/services"
	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/utils"
	"github.com/jackc/pgx/v5"
)

// CreateDoctorHandler creates a new doctor
func CreateDoctorHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.DoctorCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		doctor, err := services.NewDoctorService(db).CreateDoctorWithRequest(context.Background(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(doctor)
	}
}

// LoginDoctorHandler authenticates a doctor and returns login response
func LoginDoctorHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var loginReq models.DoctorLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		doctorService := services.NewDoctorService(db)
		doctor, err := doctorService.LoginDoctor(context.Background(), loginReq.Email, loginReq.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		token, err := utils.GenerateToken(doctor.ID, doctor.Email, "doctor")
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Create login response (without password)
		loginResp := models.DoctorLoginResponse{
			ID:         doctor.ID,
			Name:       doctor.Name,
			Email:      doctor.Email,
			Username:   doctor.Username,
			Speciality: doctor.Speciality,
			Accuracy:   doctor.Accuracy,
			Token:      token,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loginResp)
	}
}

// GetDoctorsHandler returns all doctors
func GetDoctorsHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		doctorService := services.NewDoctorService(db)
		doctors, err := doctorService.GetAllDoctors(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(doctors)
	}
}

// GetDoctorHandler returns a specific doctor
func GetDoctorHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Doctor ID is required", http.StatusBadRequest)
			return
		}

		docID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Doctor ID", http.StatusBadRequest)
			return
		}

		doctorService := services.NewDoctorService(db)
		doctor, err := doctorService.GetDoctor(context.Background(), docID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(doctor)
	}
}
