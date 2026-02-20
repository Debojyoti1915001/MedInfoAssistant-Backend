package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/models"
)

func isRetryableAIError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

// CallAIService sends image + metadata to external AI API
func CallAIService(fileBytes []byte, symptoms string, doctorSpeciality string) (*models.AIResponse, error) {

	// Convert image to base64
	base64Image := base64.StdEncoding.EncodeToString(fileBytes)

	reqBody := models.AIRequest{
		File:             base64Image,
		Symptoms:         symptoms,
		DoctorSpeciality: doctorSpeciality,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	const maxAttempts = 3
	const retryBackoff = 5 * time.Second
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		req, err := http.NewRequest(
			http.MethodPost,
			"https://rxvalidationai.onrender.com/analyze-prescription",
			bytes.NewReader(jsonData),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			if isRetryableAIError(err) && attempt < maxAttempts {
				time.Sleep(retryBackoff)
				continue
			}
			return nil, fmt.Errorf("failed to call AI service: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if (resp.StatusCode == http.StatusGatewayTimeout || resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusServiceUnavailable) && attempt < maxAttempts {
				lastErr = fmt.Errorf("AI service error: status=%d body=%s", resp.StatusCode, string(body))
				time.Sleep(retryBackoff)
				continue
			}

			return nil, fmt.Errorf("AI service error: %s", string(body))
		}

		var aiResp models.AIResponse
		if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to decode AI response: %w", err)
		}
		resp.Body.Close()

		// The AI response uses test names as map keys; copy each key into Test.Name.
		for testName, test := range aiResp.Tests {
			test.Name = testName
			aiResp.Tests[testName] = test
		}

		// The AI response uses medicine names as map keys; copy each key into Medicine.Name.
		for medicineName, medicine := range aiResp.Medicines {
			medicine.Name = medicineName
			aiResp.Medicines[medicineName] = medicine
		}

		return &aiResp, nil
	}

	if lastErr != nil {
		return nil, fmt.Errorf("failed to call AI service after retries: %w", lastErr)
	}
	return nil, fmt.Errorf("failed to call AI service after retries")
}
