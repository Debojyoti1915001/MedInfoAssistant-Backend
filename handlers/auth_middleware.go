package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/utils"
)

// AuthMiddleware wraps an http handler and verifies JWT token from Authorization header
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Verify token
		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Store claims in request context for downstream handlers
		// You can access it with: claims := r.Context().Value("claims").(*utils.Claims)
		// For now, we'll just verify the token is valid
		_ = claims

		next(w, r)
	}
}

// ExtractTokenInfo extracts user/doctor info from JWT token in request
func ExtractTokenInfo(r *http.Request) (*utils.Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("missing authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.New("invalid authorization header format")
	}

	return utils.VerifyToken(parts[1])
}
