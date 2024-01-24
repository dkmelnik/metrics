package interfaces

import "github.com/dkmelnik/metrics/internal/models"

type MetricsRepository interface {
	SaveOrUpdate(metric models.Metric) error
	SaveOrUpdateMany(metrics []models.Metric) error
	FindOneByTypeAndName(mType, mName string) (models.Metric, error)
	Find() ([]models.Metric, error)
}
