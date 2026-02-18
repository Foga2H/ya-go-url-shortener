package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"database/sql"

	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PingHandler struct {
	service *service.PingService
	logger  *logger.Logger
}

func NewPingHandler(db *sql.DB, logger *logger.Logger) *PingHandler {
	return &PingHandler{service: service.NewPingService(db, logger), logger: logger}
}

type pingResponse struct {
	Result string `json:"result"`
}

func (h *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	if err := h.service.Ping(ctx); err != nil {
		h.logger.Errorf("Failed to ping database: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(pingResponse{Result: "Success"})
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
