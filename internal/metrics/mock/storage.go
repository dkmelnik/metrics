package mock

import (
	"database/sql"
	"github.com/dkmelnik/metrics/internal/models"
	"github.com/dkmelnik/metrics/internal/storage"
)

var metricValues = map[string]map[string]interface{}{
	"gauge": {
		"HeapAlloc":  150112.000,
		"HeapSys":    3.833856e+06,
		"MCacheSys":  13.0,
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

func NewStorageMock() (*StorageMock, error) {
	r, err := storage.NewMemoryStorage("/tmp/metrics-db.json", 300, false)
	if err != nil {
		return nil, err
	}
	s := &StorageMock{r}
	s.fillStorage()
	return s, nil
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
				m.Value = sql.NullFloat64{
					Float64: fl,
					Valid:   true,
				}
			} else {
				it, _ := value.(int)
				it2 := int64(it)
				m.Delta = sql.NullInt64{
					Int64: it2,
					Valid: true,
				}
			}
			s.SaveOrUpdate(m)
		}
	}
}
