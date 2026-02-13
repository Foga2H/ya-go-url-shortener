package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/middleware"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/Foga2H/ya-go-url-shortener/internal/service"
)

type CreateLinkHandler struct {
	service *service.CreateLinkService
}

func NewCreateLinkHandler(storage repository.StorageRepo, prefixURL string) *CreateLinkHandler {
	return &CreateLinkHandler{
		service: service.NewCreateLinkService(storage, prefixURL),
	}
}

func (h *CreateLinkHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "text/plain")

	userID, ok := middleware.UserIDFromContext(req.Context())
	if !ok {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Invalid Body", http.StatusBadRequest)
		return
	}

	bodyString := string(bodyBytes)
	result, err := h.service.Create(req.Context(), userID, bodyString)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) {
			http.Error(res, "Provided URL is not valid", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrGenerateShortLink) {
			http.Error(res, "Error when trying to generate random link", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrSaveLink) {
			http.Error(res, "Error when trying to save link", http.StatusInternalServerError)
			return
		}

		http.Error(res, "Error when trying to save link", http.StatusInternalServerError)
		return
	}

	if result.IsConflict {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}
	_, err = res.Write([]byte(result.ShortURL))
	if err != nil {
		http.Error(res, "Error when trying to return response data", http.StatusBadRequest)
		return
	}
}
