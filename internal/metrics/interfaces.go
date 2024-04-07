package metrics

import (
	"context"

	"github.com/dkmelnik/metrics/internal/models"
)

// IRepository is an interface defining methods for interacting with a data repository.
type IRepository interface {
	SaveOrUpdate(ctx context.Context, metric models.Metric) error
	SaveOrUpdateMany(ctx context.Context, metrics []models.Metric) error
	FindOneByTypeAndName(ctx context.Context, mType, mName string) (models.Metric, error)
	Find(ctx context.Context) ([]models.Metric, error)
}

// Signer is an interface representing an entity capable of verifying the equality of two byte slices.
type Signer interface {
	Equal(sign, data []byte) bool
}
