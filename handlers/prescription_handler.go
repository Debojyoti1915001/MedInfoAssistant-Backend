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
	"sync"
	"time"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/models"
	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/services"
	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/utils"
	"github.com/jackc/pgx/v5"
)

type createPrescriptionResponse struct {
	Prescription *models.Prescription `json:"prescription"`
	AIAnalysis   *models.AIResponse   `json:"aiAnalysis"`
}

func persistAIItems(ctx context.Context, db *pgx.Conn, presID int64, aiResp *models.AIResponse) error {
	if aiResp == nil {
		return nil
	}

	itemService := services.NewItemsService(db)
	items := make([]*models.Items, 0, len(aiResp.Tests)+len(aiResp.Medicines))

	for testName, test := range aiResp.Tests {
		name := test.Name
		if name == "" {
			name = testName
		}

		reasonsJSON, err := json.Marshal(test)
		if err != nil {
			return fmt.Errorf("failed to encode test reasons for %s: %w", name, err)
		}

		items = append(items, &models.Items{
			PresID:    presID,
			Name:      name,
			Type:      "test",
			AIReasons: string(reasonsJSON),
			DocReason: "",
		})
	}

	for medicineName, medicine := range aiResp.Medicines {
		name := medicine.Name
		if name == "" {
			name = medicineName
		}

		reasonsJSON, err := json.Marshal(medicine)
		if err != nil {
			return fmt.Errorf("failed to encode medicine reasons for %s: %w", name, err)
		}

		items = append(items, &models.Items{
			PresID:    presID,
			Name:      name,
			Type:      "med",
			AIReasons: string(reasonsJSON),
			DocReason: "",
		})
	}

	if err := itemService.CreateItemsBulk(ctx, items); err != nil {
		return fmt.Errorf("failed to store AI items: %w", err)
	}
	return nil
}

// CreatePrescriptionHandler creates a new prescription
func CreatePrescriptionHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		const maxFileSize = 10 << 20
		const maxRequestSize = maxFileSize + (1 << 20) // file + multipart overhead

		// Enforce a hard request size limit before parsing multipart data.
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)

		// Parse multipart form with small in-memory budget; file parts spill to temp files.
		if err := r.ParseMultipartForm(2 << 20); err != nil {
			http.Error(w, "failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
			return
		}
		if r.MultipartForm != nil {
			defer r.MultipartForm.RemoveAll()
		}

		// Read file from form field 'file'
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "file is required: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		limitedReader := io.LimitReader(file, maxFileSize+1)
		fileBytes, err := io.ReadAll(limitedReader)
		if err != nil {
			http.Error(w, "failed to read file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if len(fileBytes) > maxFileSize {
			http.Error(w, "file too large (max 10MB)", http.StatusBadRequest)
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

		// Parse other form fields
		symptoms := r.FormValue("symptoms")
		userIDStr := r.FormValue("userId")
		doctorIdentifier := strings.TrimSpace(r.FormValue("doctorUsername"))

		if userIDStr == "" || doctorIdentifier == "" {
			http.Error(w, "userId and doctorUsername are required", http.StatusBadRequest)
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid userId", http.StatusBadRequest)
			return
		}

		// Resolve doctor by username or email.
		docService := services.NewDoctorService(db)
		doctor, err := docService.GetDoctorByIdentifier(context.Background(), doctorIdentifier)
		if err != nil {
			http.Error(w, "doctor not found", http.StatusNotFound)
			return
		}

		var (
			publicURL string
			aiResp    *models.AIResponse
			uploadErr error
			aiErr     error
		)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			publicURL, uploadErr = utils.UploadToSupabase(bucket, objectPath, fileBytes, contentType)
		}()
		go func() {
			defer wg.Done()
			aiResp, aiErr = services.CallAIService(fileBytes, symptoms, doctor.Speciality)
		}()
		wg.Wait()

		if uploadErr != nil {
			http.Error(w, "failed to upload file: "+uploadErr.Error(), http.StatusInternalServerError)
			return
		}
		if aiErr != nil {
			http.Error(w, "failed to analyze prescription: "+aiErr.Error(), http.StatusBadGateway)
			return
		}

		prescription := &models.Prescription{
			Symptoms: symptoms,
			Link:     publicURL,
			UserID:   userID,
			DocID:    doctor.ID,
		}

		presService := services.NewPrescriptionService(db)
		if err := presService.CreatePrescription(context.Background(), prescription); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := persistAIItems(context.Background(), db, prescription.ID, aiResp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createPrescriptionResponse{
			Prescription: prescription,
			AIAnalysis:   aiResp,
		})
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
