package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
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

	h.Storage.Set(randomString, req.URL)

	fmt.Printf("Generated link %s for %s\n", h.config.PrefixURL+"/"+randomString, req.URL)

	resp, err := json.Marshal(Response{Result: h.config.PrefixURL + "/" + randomString})
	if err != nil {
		http.Error(w, "Error when trying to marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
