package main

import (
	"net/http"

	"github.com/Foga2H/ya-go-url-shortener/internal/handler"
	storage "github.com/Foga2H/ya-go-url-shortener/internal/storage/memory"
)

func main() {
	mux := http.NewServeMux()
	memStorage := storage.NewMemStorage()

	mux.Handle(`/`, &handler.CreateLinkHandler{
		Storage: memStorage,
	})
	mux.Handle(`/{url}`, &handler.LinkHandler{
		Storage: memStorage,
	})

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
