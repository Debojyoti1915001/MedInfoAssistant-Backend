package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/services"
	"github.com/jackc/pgx/v5"
)

// AuthCheckHandler validates token and returns user/doctor info
// Frontend uses this to determine if token is valid and redirect accordingly
func AuthCheckHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract and validate token
		claims, err := ExtractTokenInfo(r)
		if err != nil {
			// No token or invalid token - user should see login/register page
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"authenticated": false,
				"message":       "No valid token. Please login/register.",
			})
			return
		}

		// Token is valid - fetch user/doctor details
		if claims.Role == "user" {
			userService := services.NewUserService(db)
			user, err := userService.GetUser(context.Background(), claims.ID)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"authenticated": false,
					"message":       "User not found",
				})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"authenticated": true,
				"role":          "user",
				"id":            user.ID,
				"name":          user.Name,
				"email":         user.Email,
				"phnNumber":     user.PhnNumber,
				"createdAt":     user.CreatedAt,
			})
		} else if claims.Role == "doctor" {
			doctorService := services.NewDoctorService(db)
			doctor, err := doctorService.GetDoctor(context.Background(), claims.ID)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"authenticated": false,
					"message":       "Doctor not found",
				})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"authenticated": true,
				"role":          "doctor",
				"id":            doctor.ID,
				"name":          doctor.Name,
				"email":         doctor.Email,
				"username":      doctor.Username,
				"speciality":    doctor.Speciality,
				"accuracy":      doctor.Accuracy,
				"phnNumber":     doctor.PhnNumber,
				"createdAt":     doctor.CreatedAt,
			})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"authenticated": false,
				"message":       "Invalid role",
			})
		}
	}
}
