package mock

import (
	"github.com/dkmelnik/metrics/internal/models"
	"github.com/dkmelnik/metrics/internal/storage"
)

var metricValues = map[string]map[string]interface{}{
	"gauge": {
		"HeapAlloc":  150112.000,
		"HeapSys":    3.833856e+06,
		"MCacheSys":  13,
		"TotalAlloc": 7.70766,
		"Mallocs":    282.00,
		"OtherSys":   3485734.100,
		"NextGC":     -3.35872e+06,
		"LastGC":     0.0,
	},
	"counter": {
		"PollCount": 14123413542,
	},
}

type StorageMock struct {
	*storage.MemoryStorage
}

func NewStorageMock() *StorageMock {
	s := &StorageMock{storage.NewMemoryStorage()}
	s.fillStorage()
	return s
}

func (s *StorageMock) fillStorage() {
	for metricType, metrics := range metricValues {
		for metricName, value := range metrics {
			var m = models.Metric{
				ID:    metricName,
				MType: metricType,
			}
			if metricType == string(models.Gauge) {
				fl, _ := value.(float64)
				m.Value = &fl
			} else {
				it, _ := value.(int64)
				m.Delta = &it
			}
			s.SaveOrUpdate(m)
		}
	}
}
