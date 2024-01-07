package interfaces

import "github.com/dkmelnik/metrics/internal/models"

type MetricsRepository interface {
	SaveOrUpdate(metric models.Metric) error
	FindOneByTypeAndID(metricType, metricID string) (int, models.Metric, error)
	GetAllMetrics() []models.Metric
}
