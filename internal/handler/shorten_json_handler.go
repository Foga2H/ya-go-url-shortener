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

type ShortenJSONHandler struct {
	service *service.ShortenService
}

func NewShortenJSONHandler(storage repository.StorageRepo, config *config.Config) *ShortenJSONHandler {
	return &ShortenJSONHandler{
		service: service.NewShortenService(storage, config.PrefixURL),
	}
}

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

func (h *ShortenJSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	jsonDecoder := json.NewDecoder(r.Body)
	var req Request
	err := jsonDecoder.Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
		return
	}

	result, err := h.service.Shorten(r.Context(), userID, req.URL)
	if err != nil {
		if errors.Is(err, service.ErrGenerateShortLink) {
			http.Error(w, "Error when trying to generate random link", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrSaveShortLink) {
			http.Error(w, "Error when trying to save link", http.StatusInternalServerError)
			return
		}

		http.Error(w, "Error when trying to save link", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(Response{Result: result.ShortURL})
	if err != nil {
		http.Error(w, "Error when trying to marshal response", http.StatusInternalServerError)
		return
	}

	if result.IsConflict {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write(resp)
}
