package metrics

import (
	"github.com/dkmelnik/metrics/internal/logger"
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
	r.Post("/update/", metricsHandler.HandleProcessMetricRequest)
	r.Get("/value/{type}/{name}", metricsHandler.HandleGetMetricValue)
	r.Post("/value/", metricsHandler.HandleGetMetric)
	r.Get("/", metricsHandler.HandleGetAllMetrics)

	return r
}
