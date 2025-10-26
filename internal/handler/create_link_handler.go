package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
)

type CreateLinkHandler struct {
	Storage repository.StorageRepo
}

func NewCreateLinkHandler(storage repository.StorageRepo) *CreateLinkHandler {
	return &CreateLinkHandler{
		Storage: storage,
	}
}

func (h *CreateLinkHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

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

	randomString, err := GenerateRandomStringURLSafe(6)
	if err != nil {
		http.Error(res, "Error when trying to generate random link", http.StatusBadRequest)
		return
	}

	h.Storage.Set(randomString, bodyString)

	fmt.Printf("Generated link %s for %s\n", randomString, bodyString)

	res.WriteHeader(http.StatusCreated)
	_, err = res.Write([]byte("http://localhost:8080/" + randomString))
	if err != nil {
		http.Error(res, "Error when trying to return response data", http.StatusBadRequest)
		return
	}
}

func GenerateRandomStringURLSafe(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
