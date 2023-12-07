package handlers

import (
	"fmt"
	"github.com/dkmelnik/metrics/internal/handlers/interfaces"
	"net/http"
)

type Handler struct {
	storage interfaces.Storage
}

func NewHandler(s interfaces.Storage) *Handler {
	return &Handler{s}
}

func (h *Handler) Create(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(rw, "Received POST request")
	}

	rw.Header().Set("content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{"status": "ok"}`))
}

func (h *Handler) GetAll(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(rw, "Received POST request")
	}
}
