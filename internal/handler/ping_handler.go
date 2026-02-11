package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"database/sql"

	"github.com/Foga2H/ya-go-url-shortener/internal/config/db"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PingHandler struct {
	*db.Config
}

func NewPingHandler(config *db.Config) *PingHandler {
	return &PingHandler{config}
}

func (h *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	dbConnection, err := sql.Open("pgx", h.Config.DatabaseDSN)
	if err != nil {
		http.Error(w, "Error when trying to connect to database", http.StatusInternalServerError)
		return
	}
	defer dbConnection.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = dbConnection.PingContext(ctx); err != nil {
		http.Error(w, "Error when trying to connect to database", http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(Response{Result: "Success"})
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
