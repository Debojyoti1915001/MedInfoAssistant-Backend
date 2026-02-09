package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/models"
	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/services"
	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/utils"
	"github.com/jackc/pgx/v5"
)

// GetUsersHandler returns all users
func GetUsersHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userService := services.NewUserService(db)
		users, err := userService.GetAllUsers(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

// CreateUserHandler creates a new user
func CreateUserHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userService := services.NewUserService(db)
		if err := userService.CreateUser(context.Background(), &user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

// LoginUserHandler authenticates a user and returns login response
func LoginUserHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var loginReq models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userService := services.NewUserService(db)
		user, err := userService.LoginUser(context.Background(), loginReq.Email, loginReq.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		token, err := utils.GenerateToken(user.ID, user.Email, "user")
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Create login response (without password)
		loginResp := models.LoginResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			PhnNumber: user.PhnNumber,
			Token:     token,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loginResp)
	}
}
