package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthHandler checks the health of the server
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "Server is running",
	})
}
