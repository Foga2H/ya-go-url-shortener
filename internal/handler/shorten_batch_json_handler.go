package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/middleware"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/Foga2H/ya-go-url-shortener/internal/service"
)

type ShortenBatchJSONHandler struct {
	service *service.ShortenBatchService
	logger  *logger.Logger
}

func NewShortenBatchJSONHandler(storage repository.StorageRepo, config *config.Config, logger *logger.Logger) *ShortenBatchJSONHandler {
	return &ShortenBatchJSONHandler{
		service: service.NewShortenBatchService(service.NewShortenService(storage, config.PrefixURL, logger)),
		logger:  logger,
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
		h.logger.Warnf("Unauthorized request: path=%s", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	jsonDecoder := json.NewDecoder(r.Body)
	var items []BatchRequest
	err := jsonDecoder.Decode(&items)
	if err != nil {
		h.logger.Errorf("Invalid JSON body: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
			h.logger.Errorf("Failed to save short links: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		h.logger.Errorf("Failed to shorten batch: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
		h.logger.Errorf("Failed to marshal response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
