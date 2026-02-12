package handler

import (
	"net/http"
	"strings"

	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

type LinkHandler struct {
	Storage repository.StorageRepo
}

func NewLinkHandler(storage repository.StorageRepo) *LinkHandler {
	return &LinkHandler{Storage: storage}
}

func (h *LinkHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "text/plain")

	id := chi.URLParam(req, "url")
	if id == "" {
		id = strings.TrimPrefix(req.URL.Path, "/")
	}

	link, ok := h.Storage.Get(req.Context(), id)
	if !ok {
		http.Error(res, "Link not found", http.StatusNotFound)
		return
	}

	res.Header().Add("Location", link)
	res.WriteHeader(http.StatusTemporaryRedirect)
	_, err := res.Write([]byte(link))
	if err != nil {
		http.Error(res, "Error when trying to return response data", http.StatusBadRequest)
		return
	}
}
