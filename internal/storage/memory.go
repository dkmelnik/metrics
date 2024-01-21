package storage

import (
	"github.com/dkmelnik/metrics/internal/apperrors"
	"github.com/dkmelnik/metrics/internal/models"
	"sync"
)

type MemoryStorage struct {
	mu          sync.RWMutex
	metrics     map[string]models.Metric
	syncsSaving bool
	filePath    string
}

func NewMemoryStorage(storagePath string, storeInterval int, restore bool) (*MemoryStorage, error) {
	ms := &MemoryStorage{
		filePath:    storagePath,
		syncsSaving: storeInterval == 0,
		metrics:     make(map[string]models.Metric),
	}

	//if restore {
	//	ms.loadMetricsFromFile()
	//}
	//
	//if storeInterval > 0 {
	//	savePeriod := time.NewTicker(time.Second * time.Duration(storeInterval))
	//	go ms.intervalUpdatingToFile(savePeriod)
	//}

	return ms, nil
}

func (m *MemoryStorage) SaveOrUpdate(metric models.Metric) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := metric.MType + "_" + metric.Name

	if existingMetric, ok := m.metrics[key]; ok {
		existingMetric.Delta = metric.Delta
		existingMetric.Value = metric.Value
		m.metrics[key] = existingMetric
	} else {
		m.metrics[key] = metric
	}

	return nil
}

func (m *MemoryStorage) FindOneByTypeAndName(mType, mName string) (models.Metric, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := mType + "_" + mName

	if metric, ok := m.metrics[key]; ok {
		return metric, nil
	}

	return models.Metric{}, apperrors.ErrNotFound
}

func (m *MemoryStorage) Find() ([]models.Metric, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var metrics []models.Metric
	for _, metric := range m.metrics {
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

//func (m *MemoryStorage) intervalUpdatingToFile(t *time.Ticker) {
//	for range t.C {
//		m.saveMetricsToFile()
//	}
//}
//
//func (m *MemoryStorage) loadMetricsFromFile() {
//	ctx := context.Background()
//
//	m.mu.Lock()
//	defer m.mu.Unlock()
//
//	file, err := os.OpenFile(m.filePath, os.O_RDONLY, 0666)
//	if err != nil {
//		logger.Log.ErrorWithContext(ctx, err)
//		return
//	}
//	defer file.Close()
//
//	var ms []models.Metric
//
//	err = json.NewDecoder(file).Decode(&m)
//	if err != nil {
//		logger.Log.ErrorWithContext(ctx, err)
//		return
//	}
//	m.metrics = make(map[string]models.Metric)
//	m.metrics = ms
//}
//
//func (m *MemoryStorage) saveMetricsToFile() {
//	ctx := context.Background()
//
//	if len(m.metrics) <= 0 {
//		return
//	}
//	m.mu.Lock()
//	defer m.mu.Unlock()
//
//	file, err := os.OpenFile(m.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
//	if err != nil {
//		logger.Log.ErrorWithContext(ctx, err)
//		return
//	}
//	defer file.Close()
//	encoder := json.NewEncoder(file)
//	if err = encoder.Encode(m.metrics); err != nil {
//		logger.Log.ErrorWithContext(ctx, err)
//		return
//	}
//}
