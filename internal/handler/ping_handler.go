package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"database/sql"

	"github.com/Foga2H/ya-go-url-shortener/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PingHandler struct {
	service *service.PingService
}

func NewPingHandler(db *sql.DB) *PingHandler {
	return &PingHandler{service: service.NewPingService(db)}
}

type pingResponse struct {
	Result string `json:"result"`
}

func (h *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	if err := h.service.Ping(ctx); err != nil {
		http.Error(w, "Error when trying to connect to database", http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(pingResponse{Result: "Success"})
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
