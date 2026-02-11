package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/models"
	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/services"
	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/utils"
	"github.com/jackc/pgx/v5"
)

// CreatePrescriptionHandler creates a new prescription
func CreatePrescriptionHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse multipart form (max 20MB in memory)
		if err := r.ParseMultipartForm(20 << 20); err != nil {
			http.Error(w, "failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Read file from form field 'file'
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "file is required: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "failed to read file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Build object path and upload
		bucket := os.Getenv("SUPABASE_STORAGE_BUCKET")
		if bucket == "" {
			bucket = "prescriptions"
		}
		// sanitize filename and add timestamp
		fname := header.Filename
		ext := filepath.Ext(fname)
		nameOnly := fname[0 : len(fname)-len(ext)]
		objectPath := fmt.Sprintf("%s_%d%s", nameOnly, time.Now().UnixNano(), ext)
		objectPath = filepath.ToSlash(objectPath)

		contentType := header.Header.Get("Content-Type")
		if contentType == "" {
			contentType = http.DetectContentType(fileBytes)
		}

		// Validate MIME type - allow only common image types
		allowedTypes := map[string]bool{
			"image/png":     true,
			"image/jpeg":    true,
			"image/jpg":     true,
			"image/gif":     true,
			"image/webp":    true,
			"image/bmp":     true,
			"image/svg+xml": true,
		}
		if !strings.HasPrefix(contentType, "image/") || !allowedTypes[contentType] {
			http.Error(w, "only image uploads are allowed", http.StatusBadRequest)
			return
		}

		// Validate extension
		extLower := strings.ToLower(ext)
		allowedExts := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true, ".bmp": true, ".svg": true}
		if !allowedExts[extLower] {
			http.Error(w, "unsupported file extension", http.StatusBadRequest)
			return
		}

		// Enforce max file size (e.g., 10MB)
		const maxSize = 10 << 20
		if len(fileBytes) > maxSize {
			http.Error(w, "file too large (max 10MB)", http.StatusBadRequest)
			return
		}

		publicURL, err := utils.UploadToSupabase(bucket, objectPath, fileBytes, contentType)
		if err != nil {
			http.Error(w, "failed to upload file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Parse other form fields
		symptoms := r.FormValue("symptoms")
		userIDStr := r.FormValue("userId")
		docIDStr := r.FormValue("docId")

		if userIDStr == "" || docIDStr == "" {
			http.Error(w, "userId and docId are required", http.StatusBadRequest)
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid userId", http.StatusBadRequest)
			return
		}

		docID, err := strconv.ParseInt(docIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid docId", http.StatusBadRequest)
			return
		}

		prescription := &models.Prescription{
			Symptoms: symptoms,
			Link:     publicURL,
			UserID:   userID,
			DocID:    docID,
		}

		presService := services.NewPrescriptionService(db)
		if err := presService.CreatePrescription(context.Background(), prescription); err != nil {
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
