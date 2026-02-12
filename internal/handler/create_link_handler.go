package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/Foga2H/ya-go-url-shortener/internal/storage/db"
	"github.com/Foga2H/ya-go-url-shortener/pkg/utils"
)

type CreateLinkHandler struct {
	Storage repository.StorageRepo
	config  *config.Config
}

func NewCreateLinkHandler(storage repository.StorageRepo, config *config.Config) *CreateLinkHandler {
	return &CreateLinkHandler{
		Storage: storage,
		config:  config,
	}
}

func (h *CreateLinkHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "text/plain")

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Invalid Body", http.StatusBadRequest)
		return
	}

	// Convert the byte slice to a string
	bodyString := string(bodyBytes)
	_, err = url.ParseRequestURI(bodyString)
	if err != nil {
		http.Error(res, "Provided URL is not valid", http.StatusBadRequest)
		return
	}

	randomString, err := utils.GenerateRandomStringURLSafe(6)
	if err != nil {
		http.Error(res, "Error when trying to generate random link", http.StatusBadRequest)
		return
	}

	storedKey, err2 := h.Storage.Set(randomString, bodyString)
	if err2 != nil {
		if errors.Is(err2, db.ErrOriginalURLConflict) {
			res.WriteHeader(http.StatusConflict)
			_, err = res.Write([]byte(h.config.PrefixURL + "/" + storedKey))
			if err != nil {
				http.Error(res, "Error when trying to return response data", http.StatusBadRequest)
				return
			}
			return
		}
		http.Error(res, "Error when trying to save link", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Generated link %s for %s\n", h.config.PrefixURL+"/"+storedKey, bodyString)

	res.WriteHeader(http.StatusCreated)
	_, err = res.Write([]byte(h.config.PrefixURL + "/" + storedKey))
	if err != nil {
		http.Error(res, "Error when trying to return response data", http.StatusBadRequest)
		return
	}
}
