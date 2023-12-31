package interfaces

type MetricsRepository interface {
	Save(metricType, metricName string, value interface{})
	FindOneByTypeName(metricType, metricName string) (interface{}, error)
	GetAllMetrics() map[string]map[string]interface{}
}
