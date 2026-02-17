package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/middleware"
	"github.com/Foga2H/ya-go-url-shortener/internal/service"
)

type DeleteUserURLsHandler struct {
	service *service.DeleteUserURLsService
}

func NewDeleteUserURLsHandler(svc *service.DeleteUserURLsService) *DeleteUserURLsHandler {
	return &DeleteUserURLsHandler{service: svc}
}

func (h *DeleteUserURLsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	jsonDecoder := json.NewDecoder(r.Body)
	var items []string
	err := jsonDecoder.Decode(&items)
	if err != nil {
		http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
		return
	}

	h.service.Enqueue(userID, items)
	w.WriteHeader(http.StatusAccepted)
}
