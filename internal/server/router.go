package server

import (
	"github.com/dkmelnik/metrics/internal/handlers"
	"github.com/dkmelnik/metrics/internal/storage"
	"net/http"
)

type User struct {
	ID        string
	FirstName string
	LastName  string
}

func (s *Server) configureRouter() {
	mux := http.NewServeMux()
	s.app.Handler = mux

	//infrastructure
	store := storage.NewCollection()
	//metrics
	metricsHandler := handlers.NewHandler(store)

	mux.HandleFunc("/update", metricsHandler.Create)
	mux.HandleFunc("/", metricsHandler.GetAll)
}
