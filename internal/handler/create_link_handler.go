package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/middleware"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/Foga2H/ya-go-url-shortener/internal/service"
)

type CreateLinkHandler struct {
	service *service.CreateLinkService
	logger  *logger.Logger
}

func NewCreateLinkHandler(storage repository.StorageRepo, prefixURL string, logger *logger.Logger) *CreateLinkHandler {
	return &CreateLinkHandler{
		service: service.NewCreateLinkService(storage, prefixURL, logger),
		logger:  logger,
	}
}

func (h *CreateLinkHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "text/plain")

	userID, ok := middleware.UserIDFromContext(req.Context())
	if !ok {
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Errorf("Failed to read body: %v", err)
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	bodyString := string(bodyBytes)
	result, err := h.service.Create(req.Context(), userID, bodyString)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) {
			h.logger.Errorf("Invalid URL: %s", req.URL.Path)
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if errors.Is(err, service.ErrGenerateShortLink) {
			h.logger.Errorf("Failed to generate short link: %v", err)
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if errors.Is(err, service.ErrSaveLink) {
			h.logger.Errorf("Error when trying to save link: %v", err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		h.logger.Errorf("Failed to create link: %v", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if result.IsConflict {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}
	_, err = res.Write([]byte(result.ShortURL))
	if err != nil {
		h.logger.Errorf("Failed to write response: %v", err)
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
