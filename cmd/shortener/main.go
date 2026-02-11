package main

import (
	"database/sql"
	"flag"
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	configDb "github.com/Foga2H/ya-go-url-shortener/internal/config/db"
	"github.com/Foga2H/ya-go-url-shortener/internal/gzip"
	"github.com/Foga2H/ya-go-url-shortener/internal/handler"
	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	dbStorage "github.com/Foga2H/ya-go-url-shortener/internal/storage/db"
	"github.com/Foga2H/ya-go-url-shortener/internal/storage/file"
	storage "github.com/Foga2H/ya-go-url-shortener/internal/storage/memory"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	r := chi.NewRouter()

	flag.Parse()
	c := config.NewConfig()
	db := configDb.NewConfig()

	var selectedStorage repository.StorageRepo = storage.NewMemStorage()

	if db.DatabaseDSN != "" {
		dbConnection, err := sql.Open("pgx", db.DatabaseDSN)
		if err != nil {
			panic(err)
		}
		selectedStorage = dbStorage.NewStorage(dbConnection)
	} else if c.FileStoragePath != "" {
		selectedStorage = file.NewStorage(c.FileStoragePath)
	}

	l := logger.NewLogger()
	gz := gzip.NewGzip()

	r.Use(l.LoggerMiddleware())
	r.Use(gz.Middleware())

	r.Post("/", handler.NewCreateLinkHandler(selectedStorage, c).ServeHTTP)
	r.Get("/{url}", handler.NewLinkHandler(selectedStorage).ServeHTTP)

	r.Get("/ping", handler.NewPingHandler(db).ServeHTTP)

	r.Post("/api/shorten", handler.NewShortenJSONHandler(selectedStorage, c).ServeHTTP)

	err := http.ListenAndServe(c.BaseURL, r)
	if err != nil {
		panic(err)
	}
}
