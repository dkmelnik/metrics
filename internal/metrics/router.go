package metrics

import (
	"github.com/dkmelnik/metrics/configs"
	"github.com/dkmelnik/metrics/internal/metrics/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/dkmelnik/metrics/internal/middlewares"
	"github.com/dkmelnik/metrics/internal/storage"
)

func ConfigureRouter(pgDB *sqlx.DB, storageConfig configs.Server) (*chi.Mux, error) {
	r := chi.NewRouter()

	//r.Use(logger.Log.RequestLog)

	//infrastructure

	var store interfaces.MetricsRepository
	var err error
	if pgDB == nil {
		store, err = storage.NewMemoryStorage(storageConfig.FileStoragePath, storageConfig.StoreInterval, storageConfig.Restore)
	} else {
		store, err = storage.NewRepositoryStorage(pgDB)
	}
	if err != nil {
		return nil, err
	}

	//metrics
	service := NewService(store)
	metricsHandler := NewHandler(pgDB, service)

	r.Post("/update/{type}/{name}/{value}", metricsHandler.CreateOrUpdateByParams)

	r.Route("/update/", func(r chi.Router) {
		r.Use(middlewares.Compress)
		r.Post("/", metricsHandler.CreateOrUpdateByJSON)
	})

	r.Get("/value/{type}/{name}", metricsHandler.GetMetricValue)

	r.Route("/value/", func(r chi.Router) {
		r.Use(middlewares.Compress)
		r.Post("/", metricsHandler.GetMetric)
	})

	r.Route("/", func(r chi.Router) {
		r.Use(middlewares.Compress)
		r.Get("/", metricsHandler.GetAllMetrics)
	})

	r.Get("/ping", metricsHandler.CheckPostgresDBConnection)

	return r, nil
}
