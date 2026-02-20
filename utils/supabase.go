package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

// UploadToSupabase uploads the provided file bytes to the given Supabase Storage bucket and path.
// It returns the public URL for the uploaded object (assumes the bucket is public), or an error.
func UploadToSupabase(bucket, objectPath string, fileBytes []byte, contentType string) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	if supabaseKey == "" {
		// Backward compatibility with existing env var name.
		supabaseKey = os.Getenv("SUPABASE_SERVICE_KEY")
	}
	if supabaseURL == "" || supabaseKey == "" {
		return "", fmt.Errorf("SUPABASE_URL or SUPABASE_SERVICE_ROLE_KEY not set")
	}

	// Supabase Storage Authorization expects a JWT (service role / anon legacy key).
	// Keys like sb_publishable_* are not JWTs and will fail with "Invalid Compact JWS".
	if strings.HasPrefix(supabaseKey, "sb_publishable_") {
		return "", fmt.Errorf("invalid Supabase key for server upload: use SUPABASE_SERVICE_ROLE_KEY (JWT), not publishable key")
	}

	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, bucket, objectPath)

	req, err := http.NewRequest(http.MethodPut, uploadURL, bytes.NewReader(fileBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", supabaseKey))
	req.Header.Set("apikey", supabaseKey)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// Use no-cache to avoid caching
	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	// Construct public URL - assumes object is in a public bucket
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, bucket, objectPath)
	return publicURL, nil
}

// Helper to sanitize file name and append timestamp or unique suffix if needed
func BuildObjectPath(folder, filename string) string {
	base := path.Base(filename)
	return path.Join(folder, base)
}
