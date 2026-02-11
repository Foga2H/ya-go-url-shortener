package main

import (
	"flag"
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	"github.com/Foga2H/ya-go-url-shortener/internal/gzip"
	"github.com/Foga2H/ya-go-url-shortener/internal/handler"
	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/storage/file"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	flag.Parse()
	c := config.NewConfig()
	//memStorage := storage.NewMemStorage()
	fileStorage := file.NewStorage(c.FileStoragePath)
	l := logger.NewLogger()
	gz := gzip.NewGzip()

	r.Use(l.LoggerMiddleware())
	r.Use(gz.Middleware())

	r.Post("/", handler.NewCreateLinkHandler(fileStorage, c).ServeHTTP)
	r.Get("/{url}", handler.NewLinkHandler(fileStorage).ServeHTTP)

	r.Post("/api/shorten", handler.NewShortenJSONHandler(fileStorage, c).ServeHTTP)

	err := http.ListenAndServe(c.BaseURL, r)
	if err != nil {
		panic(err)
	}
}
