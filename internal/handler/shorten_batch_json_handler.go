package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/Foga2H/ya-go-url-shortener/internal/storage/db"
	"github.com/Foga2H/ya-go-url-shortener/pkg/utils"
)

type ShortenBatchJSONHandler struct {
	Storage repository.StorageRepo
	config  *config.Config
}

func NewShortenBatchJSONHandler(storage repository.StorageRepo, config *config.Config) *ShortenBatchJSONHandler {
	return &ShortenBatchJSONHandler{
		Storage: storage,
		config:  config,
	}
}

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResult struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (h *ShortenBatchJSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	jsonDecoder := json.NewDecoder(r.Body)
	var items []BatchRequest
	err := jsonDecoder.Decode(&items)
	if err != nil {
		http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
		return
	}

	var results []BatchResult

	for _, item := range items {
		result, err := h.shortItem(item)
		if err != nil {
			http.Error(w, "Error when trying to save link", http.StatusInternalServerError)
			return
		}
		results = append(results, result)
	}

	resp, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Error when trying to marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *ShortenBatchJSONHandler) shortItem(item BatchRequest) (BatchResult, error) {
	randomString, err := utils.GenerateRandomStringURLSafe(6)
	if err != nil {
		return BatchResult{}, err
	}

	storedKey, err := h.Storage.Set(randomString, item.OriginalURL)
	if err != nil {
		if errors.Is(err, db.ErrOriginalURLConflict) {
			var result BatchResult
			result.CorrelationID = item.CorrelationID
			result.ShortURL = h.config.PrefixURL + "/" + storedKey
			return result, nil
		}
		return BatchResult{}, err
	}
	fmt.Printf("Generated link %s for %s\n", h.config.PrefixURL+"/"+storedKey, item.OriginalURL)

	var result BatchResult
	result.CorrelationID = item.CorrelationID
	result.ShortURL = h.config.PrefixURL + "/" + storedKey

	return result, nil
}
