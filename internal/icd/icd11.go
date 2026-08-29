package icd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type DestinationEntity struct {
	Title string  `json:"title"`
	Code  string  `json:"theCode"`
	Score float64 `json:"score"`
}

type ICD11SearchResponse struct {
	DestinationEntities []DestinationEntity `json:"destinationEntities"`
	Error               bool                `json:"error,omitempty"`
	ErrorMessage        string              `json:"errorMessage,omitempty"`
}

// SearchICD11 searches for the disease with the given name and
// returns entities mathcing the disease.
// It uses the ICD-API with the ICD 11 codes.
func SearchICD11(diseaseName string) ([]DestinationEntity, error) {
	// ICD-API url for searching for codes in the mms linearization for the
	// 2026-01 release
	baseURL := "http://localhost:8000/icd/release/11/2026-01/mms/search"

	// Query parameters
	params := url.Values{}
	params.Add("q", diseaseName)
	params.Add("highlightingEnabled", "false") // Disable highlighting to remove em tags

	// Build full url
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Create the request
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("API-Version", "v2")
	req.Header.Add("Accept-Language", "en")

	client := http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get: %w", err)
	}
	defer resp.Body.Close()

	// Check for error responses (non-200)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("icd-10 api returned status %d: %s", resp.StatusCode, string(body))
	}

	var searchResults ICD11SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResults); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for any errors in the results
	if searchResults.Error {
		return nil, fmt.Errorf("error searching codes: %s", searchResults.ErrorMessage)
	}

	filteredEntities := filterDestinationEntities(searchResults.DestinationEntities, -2)

	return filteredEntities, nil
}

// filterDestinationEntities removes entities with a lower score than the minScore
func filterDestinationEntities(
	entities []DestinationEntity,
	minScore float64,
) (newEntities []DestinationEntity) {
	for _, entity := range entities {
		if entity.Score > minScore {
			newEntities = append(newEntities, entity)
		}
	}
	return
}
