package metrics

import (
	"context"
	"github.com/dkmelnik/metrics/internal/models"
)

type Repository interface {
	SaveOrUpdate(ctx context.Context, metric models.Metric) error
	SaveOrUpdateMany(ctx context.Context, metrics []models.Metric) error
	FindOneByTypeAndName(ctx context.Context, mType, mName string) (models.Metric, error)
	Find(ctx context.Context) ([]models.Metric, error)
}
