package metrics

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/dkmelnik/metrics/internal/middlewares"
)

func ConfigureRouter(
	pgDB *sqlx.DB,
	storage Repository,
	signer middlewares.Signer,
) (*chi.Mux, error) {
	r := chi.NewRouter()

	//r.Use(logger.Log.RequestLog)

	// infrastructure
	m := middlewares.NewMiddlewareManager(signer)
	r.Use(m.Recovery)

	//metrics
	service := NewService(storage)
	metricsHandler := NewHandler(pgDB, service)

	r.Post("/update/{type}/{name}/{value}", metricsHandler.CreateOrUpdateByParams)

	r.Route("/update/", func(r chi.Router) {
		r.Use(m.Sign)
		r.Use(m.Compress)
		r.Post("/", metricsHandler.CreateOrUpdateByJSON)
	})

	r.Route("/updates/", func(r chi.Router) {
		r.Use(m.Compress)
		r.Post("/", metricsHandler.CreateOrUpdateMany)
	})

	r.Get("/value/{type}/{name}", metricsHandler.GetMetricValue)

	r.Route("/value/", func(r chi.Router) {
		r.Use(m.Compress)
		r.Post("/", metricsHandler.GetMetric)
	})

	r.Route("/", func(r chi.Router) {
		r.Use(m.Compress)
		r.Get("/", metricsHandler.GetAllMetrics)
	})

	r.Get("/ping", metricsHandler.CheckPostgresDBConnection)

	return r, nil
}
