package server

import (
	"github.com/dkmelnik/metrics/internal/handlers"
	"github.com/dkmelnik/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

type User struct {
	ID        string
	FirstName string
	LastName  string
}

func (s *Server) configureRouter() {
	r := chi.NewRouter()

	s.app.Handler = r

	//infrastructure
	store := storage.NewMemoryStorage()
	//metrics
	metricsHandler := handlers.NewHandler(store)

	r.Post("/update/{type}/{name}/{value}", metricsHandler.Create)
	r.Get("/value/{type}/{name}", metricsHandler.Get)
	r.Get("/", metricsHandler.GetAll)
}
