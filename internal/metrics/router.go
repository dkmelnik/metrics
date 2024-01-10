package metrics

import (
	"github.com/dkmelnik/metrics/configs"
	"github.com/go-chi/chi/v5"

	"github.com/dkmelnik/metrics/internal/logger"
	"github.com/dkmelnik/metrics/internal/middlewares"
	"github.com/dkmelnik/metrics/internal/storage"
)

func ConfigureRouter(c configs.Storage) (*chi.Mux, error) {
	r := chi.NewRouter()

	r.Use(logger.Log.RequestLogger)

	//infrastructure
	store, err := storage.NewMemoryStorage(c.FileStoragePath, c.StoreInterval, c.Restore)
	if err != nil {
		return nil, err
	}
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

	return r, nil
}
