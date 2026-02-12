package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	"github.com/Foga2H/ya-go-url-shortener/internal/middleware"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/Foga2H/ya-go-url-shortener/internal/storage/db"
	"github.com/Foga2H/ya-go-url-shortener/pkg/utils"
)

type ShortenJSONHandler struct {
	Storage repository.StorageRepo
	config  *config.Config
}

func NewShortenJSONHandler(storage repository.StorageRepo, config *config.Config) *ShortenJSONHandler {
	return &ShortenJSONHandler{
		Storage: storage,
		config:  config,
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

	randomString, err := utils.GenerateRandomStringURLSafe(6)
	if err != nil {
		http.Error(w, "Error when trying to generate random link", http.StatusBadRequest)
		return
	}

	storedKey, err := h.Storage.Set(r.Context(), userID, randomString, req.URL)
	if err != nil {
		if errors.Is(err, db.ErrOriginalURLConflict) {
			w.WriteHeader(http.StatusConflict)
			resp, err := json.Marshal(Response{Result: h.config.PrefixURL + "/" + storedKey})
			if err != nil {
				http.Error(w, "Error when trying to marshal response", http.StatusInternalServerError)
				return
			}
			w.Write(resp)
			return
		}
		http.Error(w, "Error when trying to save link", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Generated link %s for %s\n", h.config.PrefixURL+"/"+storedKey, req.URL)

	resp, err := json.Marshal(Response{Result: h.config.PrefixURL + "/" + storedKey})
	if err != nil {
		http.Error(w, "Error when trying to marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
