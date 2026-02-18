package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/auth"
	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	configDb "github.com/Foga2H/ya-go-url-shortener/internal/config/db"
	"github.com/Foga2H/ya-go-url-shortener/internal/gzip"
	"github.com/Foga2H/ya-go-url-shortener/internal/handler"
	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/middleware"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/Foga2H/ya-go-url-shortener/internal/service"
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

	l := logger.NewLogger()

	if db.DatabaseDSN != "" {
		dbConnection, err := sql.Open("pgx", db.DatabaseDSN)
		if err != nil {
			l.Fatalf("failed to connect to postgres: %v", err)
		}

		driver, err := postgres.WithInstance(dbConnection, &postgres.Config{})
		if err != nil {
			log.Fatalf("failed to connect to postgres: %v", err)
		}

		m, err := migrate.NewWithDatabaseInstance(
			"file://migrations",
			"postgres", driver)
		if err != nil {
			log.Fatalf("failed to connect to postgres: %v", err)
		}

		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("failed to run migrations: %v", err)
		}

		pingDB = dbConnection
		selectedStorage = dbStorage.NewStorage(dbConnection)
	} else if c.FileStoragePath != "" {
		selectedStorage = file.NewStorage(c.FileStoragePath)
	}

	gz := gzip.NewGzip()
	token := auth.NewToken(c.JWTSecret)
	uc := middleware.NewUserCookie(token)

	r.Use(l.LoggerMiddleware())
	r.Use(gz.Middleware())
	r.Use(uc.Middleware())

	if pingDB != nil {
		r.Get("/ping", handler.NewPingHandler(pingDB, l).ServeHTTP)
	}

	r.Post("/", handler.NewCreateLinkHandler(selectedStorage, c.PrefixURL, l).ServeHTTP)
	r.Get("/{url}", handler.NewLinkHandler(selectedStorage, l).ServeHTTP)

	r.Post("/api/shorten", handler.NewShortenJSONHandler(selectedStorage, c, l).ServeHTTP)
	r.Post("/api/shorten/batch", handler.NewShortenBatchJSONHandler(selectedStorage, c, l).ServeHTTP)
	r.Get("/api/user/urls", handler.NewUserURLsHandler(selectedStorage, c, l).ServeHTTP)

	deleteSvc := service.NewDeleteUserURLsService(selectedStorage, c, l)
	deleteSvc.Start(context.Background())

	r.Delete("/api/user/urls", handler.NewDeleteUserURLsHandler(deleteSvc, l).ServeHTTP)

	err := http.ListenAndServe(c.BaseURL, r)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
