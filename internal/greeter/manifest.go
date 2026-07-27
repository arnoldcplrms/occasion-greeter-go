package greeter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func GetManifestFromURL(url string) (*OccasionManifest, error) {
	normalized := NormalizeManifestSourceURL(url)
	if normalized == "" {
		return nil, fmt.Errorf("OCCASION_PHOTO_MANIFEST_URL environment variable is not set")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(normalized)
	if err != nil {
		return nil, fmt.Errorf("fetch manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("manifest request failed (status %d)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}

	var manifest OccasionManifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}

	return &manifest, nil
}

func GetManifest() (*OccasionManifest, error) {
	url := os.Getenv("OCCASION_PHOTO_MANIFEST_URL")
	return GetManifestFromURL(url)
}
