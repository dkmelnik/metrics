package grpc

import (
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"

	servicemetrics "github.com/dkmelnik/metrics/internal/metrics"
	grpcmetrics "github.com/dkmelnik/metrics/proto/metrics"
)

func ConfigureRouter(
	app *grpc.Server,
	pgDB *sqlx.DB,
	storage servicemetrics.IRepository,
) error {
	service := servicemetrics.NewService(storage)
	metricsHandler := NewHandler(pgDB, service)
	grpcmetrics.RegisterMetricsServer(app, metricsHandler)
	return nil
}
