package main

import (
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/handler"
	storage "github.com/Foga2H/ya-go-url-shortener/internal/storage/memory"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	memStorage := storage.NewMemStorage()

	r.Post("/", handler.NewLinkHandler(memStorage).ServeHTTP)
	r.Get("/{url}", handler.NewCreateLinkHandler(memStorage).ServeHTTP)

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}
