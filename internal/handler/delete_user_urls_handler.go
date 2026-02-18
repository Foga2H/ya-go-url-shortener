package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/middleware"
	"github.com/Foga2H/ya-go-url-shortener/internal/service"
)

type DeleteUserURLsHandler struct {
	service *service.DeleteUserURLsService
	logger  *logger.Logger
}

func NewDeleteUserURLsHandler(svc *service.DeleteUserURLsService, logger *logger.Logger) *DeleteUserURLsHandler {
	return &DeleteUserURLsHandler{service: svc, logger: logger}
}

func (h *DeleteUserURLsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	jsonDecoder := json.NewDecoder(r.Body)
	var items []string
	err := jsonDecoder.Decode(&items)
	if err != nil {
		h.logger.Errorf("Failed to decode body: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.service.Enqueue(userID, items)
	w.WriteHeader(http.StatusAccepted)
}
