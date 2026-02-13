package main

import (
	"database/sql"
	"errors"
	"flag"
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/auth"
	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	configDb "github.com/Foga2H/ya-go-url-shortener/internal/config/db"
	"github.com/Foga2H/ya-go-url-shortener/internal/gzip"
	"github.com/Foga2H/ya-go-url-shortener/internal/handler"
	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/middleware"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	dbStorage "github.com/Foga2H/ya-go-url-shortener/internal/storage/db"
	"github.com/Foga2H/ya-go-url-shortener/internal/storage/file"
	storage "github.com/Foga2H/ya-go-url-shortener/internal/storage/memory"
	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	r := chi.NewRouter()

	flag.Parse()
	c := config.NewConfig()
	db := configDb.NewConfig()

	var selectedStorage repository.StorageRepo = storage.NewMemStorage()
	var pingDB *sql.DB

	if db.DatabaseDSN != "" {
		dbConnection, err := sql.Open("pgx", db.DatabaseDSN)
		if err != nil {
			panic(err)
		}

		driver, err := postgres.WithInstance(dbConnection, &postgres.Config{})
		if err != nil {
			panic(err)
		}

		m, err := migrate.NewWithDatabaseInstance(
			"file://migrations",
			"postgres", driver)
		if err != nil {
			panic(err)
		}

		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			panic(err)
		}

		pingDB = dbConnection
		selectedStorage = dbStorage.NewStorage(dbConnection)
	} else if c.FileStoragePath != "" {
		selectedStorage = file.NewStorage(c.FileStoragePath)
	}

	l := logger.NewLogger()
	gz := gzip.NewGzip()
	token := auth.NewToken(c.JWTSecret)
	uc := middleware.NewUserCookie(token)

	r.Use(l.LoggerMiddleware())
	r.Use(gz.Middleware())
	r.Use(uc.Middleware())

	if pingDB != nil {
		r.Get("/ping", handler.NewPingHandler(pingDB).ServeHTTP)
	}

	r.Post("/", handler.NewCreateLinkHandler(selectedStorage, c.PrefixURL).ServeHTTP)
	r.Get("/{url}", handler.NewLinkHandler(selectedStorage).ServeHTTP)

	r.Post("/api/shorten", handler.NewShortenJSONHandler(selectedStorage, c).ServeHTTP)
	r.Post("/api/shorten/batch", handler.NewShortenBatchJSONHandler(selectedStorage, c).ServeHTTP)
	r.Get("/api/user/urls", handler.NewUserURLsHandler(selectedStorage, c).ServeHTTP)

	err := http.ListenAndServe(c.BaseURL, r)
	if err != nil {
		panic(err)
	}
}
