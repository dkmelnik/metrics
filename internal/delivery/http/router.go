package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/dkmelnik/metrics/internal/delivery/http/metrics"
	"github.com/dkmelnik/metrics/internal/logger"
	storage "github.com/dkmelnik/metrics/internal/metrics"
	"github.com/dkmelnik/metrics/internal/middlewares/http"
)

func ConfigureRouter(
	trustedSubnet string,
	privateKeyPath string,
	pgDB *sqlx.DB,
	repo storage.IRepository,
	signer http.Signer,
) (*chi.Mux, error) {
	r := chi.NewRouter()

	r.Use(logger.Log.RequestLog)

	// infrastructure
	m, err := http.NewMiddlewareManager(trustedSubnet, privateKeyPath, signer)
	if err != nil {
		return nil, err
	}
	r.Use(m.Recovery)
	r.Use(m.TrustedSubnet)

	//metrics
	service := storage.NewService(repo)
	metricsHandler := metrics.NewHandler(pgDB, service)

	r.Post("/update/{type}/{name}/{value}", metricsHandler.CreateOrUpdateByParams)

	r.Route("/update/", func(r chi.Router) {
		r.Use(m.Sign)
		r.Use(m.Decryption)
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
