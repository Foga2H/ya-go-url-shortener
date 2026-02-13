package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	"github.com/Foga2H/ya-go-url-shortener/internal/middleware"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/Foga2H/ya-go-url-shortener/internal/service"
)

type ShortenBatchJSONHandler struct {
	service *service.ShortenBatchService
}

func NewShortenBatchJSONHandler(storage repository.StorageRepo, config *config.Config) *ShortenBatchJSONHandler {
	return &ShortenBatchJSONHandler{
		service: service.NewShortenBatchService(service.NewShortenService(storage, config.PrefixURL)),
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

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	jsonDecoder := json.NewDecoder(r.Body)
	var items []BatchRequest
	err := jsonDecoder.Decode(&items)
	if err != nil {
		http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
		return
	}

	serviceItems := make([]service.BatchItem, 0, len(items))
	for _, item := range items {
		serviceItems = append(serviceItems, service.BatchItem{
			CorrelationID: item.CorrelationID,
			OriginalURL:   item.OriginalURL,
		})
	}

	results, err := h.service.Shorten(r.Context(), userID, serviceItems)
	if err != nil {
		if errors.Is(err, service.ErrGenerateShortLink) || errors.Is(err, service.ErrSaveShortLink) {
			http.Error(w, "Error when trying to save link", http.StatusInternalServerError)
			return
		}
		http.Error(w, "Error when trying to save link", http.StatusInternalServerError)
		return
	}

	respResult := make([]BatchResult, 0, len(results))
	for _, item := range results {
		respResult = append(respResult, BatchResult{
			CorrelationID: item.CorrelationID,
			ShortURL:      item.ShortURL,
		})
	}

	resp, err := json.Marshal(respResult)
	if err != nil {
		http.Error(w, "Error when trying to marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
