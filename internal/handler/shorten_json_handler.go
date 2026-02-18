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

type ShortenJSONHandler struct {
	service *service.ShortenService
	logger  *logger.Logger
}

func NewShortenJSONHandler(storage repository.StorageRepo, config *config.Config, logger *logger.Logger) *ShortenJSONHandler {
	return &ShortenJSONHandler{
		service: service.NewShortenService(storage, config.PrefixURL, logger),
		logger:  logger,
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
		h.logger.Warnf("Unauthorized request: path=%s", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	jsonDecoder := json.NewDecoder(r.Body)
	var req Request
	err := jsonDecoder.Decode(&req)
	if err != nil {
		h.logger.Errorf("Invalid JSON body: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := h.service.Shorten(r.Context(), userID, req.URL)
	if err != nil {
		if errors.Is(err, service.ErrGenerateShortLink) {
			h.logger.Errorf("Failed to generate short link: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrSaveShortLink) {
			h.logger.Errorf("Failed to save short link: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		h.logger.Errorf("Failed to shorten URL: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(Response{Result: result.ShortURL})
	if err != nil {
		h.logger.Errorf("Failed to marshal response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if result.IsConflict {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write(resp)
}
