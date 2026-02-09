package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/services"
	"github.com/jackc/pgx/v5"
)

// UserProfileHandler returns authenticated user's profile
func UserProfileHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract token claims
		claims, err := ExtractTokenInfo(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Verify user role
		if claims.Role != "user" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Fetch user profile
		userService := services.NewUserService(db)
		user, err := userService.GetUser(context.Background(), claims.ID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":        user.ID,
			"name":      user.Name,
			"email":     user.Email,
			"phnNumber": user.PhnNumber,
			"createdAt": user.CreatedAt,
		})
	}
}

// DoctorProfileHandler returns authenticated doctor's profile
func DoctorProfileHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract token claims
		claims, err := ExtractTokenInfo(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Verify doctor role
		if claims.Role != "doctor" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Fetch doctor profile
		doctorService := services.NewDoctorService(db)
		doctor, err := doctorService.GetDoctor(context.Background(), claims.ID)
		if err != nil {
			http.Error(w, "Doctor not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":         doctor.ID,
			"name":       doctor.Name,
			"email":      doctor.Email,
			"username":   doctor.Username,
			"speciality": doctor.Speciality,
			"accuracy":   doctor.Accuracy,
			"phnNumber":  doctor.PhnNumber,
			"createdAt":  doctor.CreatedAt,
		})
	}
}
