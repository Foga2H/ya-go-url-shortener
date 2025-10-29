package main

import (
	"flag"
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/config"
	"github.com/Foga2H/ya-go-url-shortener/internal/handler"
	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	storage "github.com/Foga2H/ya-go-url-shortener/internal/storage/memory"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	flag.Parse()
	c := config.NewConfig()
	memStorage := storage.NewMemStorage()
	l := logger.NewLogger()

	r.Use(l.LoggerMiddleware())

	r.Post("/", handler.NewCreateLinkHandler(memStorage, c).ServeHTTP)
	r.Get("/{url}", handler.NewLinkHandler(memStorage).ServeHTTP)

	r.Post("/api/shorten", handler.NewShortenJSONHandler(memStorage, c).ServeHTTP)

	err := http.ListenAndServe(c.BaseURL, r)
	if err != nil {
		panic(err)
	}
}
