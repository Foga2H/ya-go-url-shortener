package main

import (
	"net/http"

	config "github.com/Foga2H/ya-go-url-shortener/internal/config"
	"github.com/Foga2H/ya-go-url-shortener/internal/handler"
	storage "github.com/Foga2H/ya-go-url-shortener/internal/storage/memory"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	c := config.NewConfig()
	memStorage := storage.NewMemStorage()

	r.Post("/", handler.NewCreateLinkHandler(memStorage, c).ServeHTTP)
	r.Get("/{url}", handler.NewLinkHandler(memStorage).ServeHTTP)

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}
