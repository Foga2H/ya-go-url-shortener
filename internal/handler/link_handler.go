package handler

import (
	"net/http"
	"strings"

	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
)

type LinkHandler struct {
	Storage repository.StorageRepo
}

func NewLinkHandler(storage repository.StorageRepo) *LinkHandler {
	return &LinkHandler{}
}

func (h *LinkHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	res.Header().Set("content-type", "text/plain")

	id := req.PathValue("url")
	if id == "" {
		id = strings.TrimPrefix(req.URL.Path, "/")
	}

	link, ok := h.Storage.Get(id)
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
