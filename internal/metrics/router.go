package metrics

import (
	"github.com/go-chi/chi/v5"

	"github.com/dkmelnik/metrics/internal/storage"
)

func ConfigureRouter() *chi.Mux {
	r := chi.NewRouter()

	//infrastructure
	store := storage.NewMemoryStorage()

	//metrics
	service := NewService(store)
	metricsHandler := NewHandler(service)

	r.Post("/update/{type}/{name}/{value}", metricsHandler.Create)
	r.Get("/value/{type}/{name}", metricsHandler.Get)
	r.Get("/", metricsHandler.GetAll)

	return r
}
