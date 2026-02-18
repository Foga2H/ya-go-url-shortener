package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/middleware"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/Foga2H/ya-go-url-shortener/internal/service"
)

type UserURLsHandler struct {
	service *service.UserURLsService
	logger  *logger.Logger
}

type userURLResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewUserURLsHandler(storage repository.StorageRepo, config *config.Config, logger *logger.Logger) *UserURLsHandler {
	return &UserURLsHandler{
		service: service.NewUserURLsService(storage, config.PrefixURL, logger),
		logger:  logger,
	}
}

func (h *UserURLsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		h.logger.Warnf("Unauthorized request: path=%s", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	results, err := h.service.List(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to list user URLs: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(results) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response := make([]userURLResponse, 0, len(results))
	for _, item := range results {
		response = append(response, userURLResponse{
			ShortURL:    item.ShortURL,
			OriginalURL: item.OriginalURL,
		})
	}

	resp, err := json.Marshal(response)
	if err != nil {
		h.logger.Errorf("Failed to marshal response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
