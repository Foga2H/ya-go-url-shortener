package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	"github.com/Foga2H/ya-go-url-shortener/internal/middleware"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
)

type UserURLsHandler struct {
	storage repository.StorageRepo
	config  *config.Config
}

type userURLResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewUserURLsHandler(storage repository.StorageRepo, config *config.Config) *UserURLsHandler {
	return &UserURLsHandler{
		storage: storage,
		config:  config,
	}
}

func (h *UserURLsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	results, err := h.storage.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	if len(results) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	prefix := strings.TrimRight(h.config.PrefixURL, "/") + "/"
	response := make([]userURLResponse, 0, len(results))
	for _, item := range results {
		response = append(response, userURLResponse{
			ShortURL:    prefix + strings.TrimLeft(item.ShortURL, "/"),
			OriginalURL: item.OriginalURL,
		})
	}

	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error when trying to marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
