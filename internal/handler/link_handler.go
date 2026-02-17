package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/Foga2H/ya-go-url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type LinkHandler struct {
	service *service.LinkService
}

func NewLinkHandler(storage repository.StorageRepo) *LinkHandler {
	return &LinkHandler{service: service.NewLinkService(storage)}
}

func (h *LinkHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "text/plain")

	id := chi.URLParam(req, "url")
	if id == "" {
		id = strings.TrimPrefix(req.URL.Path, "/")
	}

	link, err := h.service.Resolve(req.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrLinkNotFound) {
			http.Error(res, "Link not found", http.StatusNotFound)
			return
		}
		
		if errors.Is(err, service.ErrLinkIsDeleted) {
			http.Error(res, "Link is deleted", http.StatusGone)
			return
		}

		http.Error(res, "Link not found", http.StatusNotFound)
		return
	}

	res.Header().Add("Location", link)
	res.WriteHeader(http.StatusTemporaryRedirect)
	_, err = res.Write([]byte(link))
	if err != nil {
		http.Error(res, "Error when trying to return response data", http.StatusBadRequest)
		return
	}
}
