package metrics

import (
	"github.com/dkmelnik/metrics/internal/logger"
	"github.com/dkmelnik/metrics/internal/middlewares"
	"github.com/go-chi/chi/v5"

	"github.com/dkmelnik/metrics/internal/storage"
)

func ConfigureRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(logger.Log.RequestLogger)

	//infrastructure
	store := storage.NewMemoryStorage()

	//metrics
	service := NewService(store)
	metricsHandler := NewHandler(service)

	r.Post("/update/{type}/{name}/{value}", metricsHandler.HandleRecordMetricValue)

	r.Route("/update/", func(r chi.Router) {
		r.Use(middlewares.Compress)
		r.Post("/", metricsHandler.HandleProcessMetricRequest)
	})

	r.Get("/value/{type}/{name}", metricsHandler.HandleGetMetricValue)

	r.Route("/value/", func(r chi.Router) {
		r.Use(middlewares.Compress)
		r.Post("/", metricsHandler.HandleGetMetric)
	})

	r.Route("/", func(r chi.Router) {
		r.Use(middlewares.Compress)
		r.Get("/", metricsHandler.HandleGetAllMetrics)
	})

	return r
}
