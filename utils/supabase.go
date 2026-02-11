package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

// UploadToSupabase uploads the provided file bytes to the given Supabase Storage bucket and path.
// It returns the public URL for the uploaded object (assumes the bucket is public), or an error.
func UploadToSupabase(bucket, objectPath string, fileBytes []byte, contentType string) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_KEY")
	if supabaseURL == "" || supabaseKey == "" {
		return "", fmt.Errorf("SUPABASE_URL or SUPABASE_SERVICE_KEY not set")
	}

	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, bucket, objectPath)

	req, err := http.NewRequest(http.MethodPut, uploadURL, bytes.NewReader(fileBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", supabaseKey))
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