package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Foga2H/ya-go-url-shortener/internal/logger"
	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/Foga2H/ya-go-url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type LinkHandler struct {
	service *service.LinkService
	logger  *logger.Logger
}

func NewLinkHandler(storage repository.StorageRepo, logger *logger.Logger) *LinkHandler {
	return &LinkHandler{service: service.NewLinkService(storage, logger), logger: logger}
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
			h.logger.Warnf("Link not found: id=%s", id)
			http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if errors.Is(err, service.ErrLinkIsDeleted) {
			h.logger.Warnf("Link is deleted: id=%s", id)
			http.Error(res, http.StatusText(http.StatusGone), http.StatusGone)
			return
		}

		h.logger.Errorf("Failed to resolve link id=%s: %v", id, err)
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	res.Header().Add("Location", link)
	res.WriteHeader(http.StatusTemporaryRedirect)
	_, err = res.Write([]byte(link))
	if err != nil {
		h.logger.Errorf("Failed to write response: %v", err)
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
