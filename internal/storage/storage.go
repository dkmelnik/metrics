package storage

import (
	"errors"
	"github.com/dkmelnik/metrics/internal/models"
	"strings"
	"sync"
)

type MemoryStorage struct {
	mu      sync.RWMutex
	metrics []models.Metric
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		metrics: make([]models.Metric, 0),
	}
}

func (ms *MemoryStorage) SaveOrUpdate(metric models.Metric) error {
	if idx, prev, err := ms.FindOneByTypeAndID(metric.MType, metric.ID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			ms.mu.Lock()
			ms.metrics = append(ms.metrics, metric)
			ms.mu.Unlock()
			return nil
		} else {
			return err
		}
	} else {
		prev.Value = metric.Value
		prev.Delta = metric.Delta
		if err = ms.UpdateByIDX(idx, prev); err != nil {
			return err
		}
	}
	return nil
}

func (ms *MemoryStorage) UpdateByIDX(idx int, newValue models.Metric) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if len(ms.metrics) < idx {
		return errors.New("metric not found")
	}

	ms.metrics[idx] = newValue

	return nil
}

func (ms *MemoryStorage) FindOneByTypeAndID(metricType, metricID string) (int, models.Metric, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	for idx, metric := range ms.metrics {
		if metric.MType == metricType && metric.ID == metricID {
			return idx, metric, nil
		}
	}
	return 0, models.Metric{}, errors.New("metric not found")
}

func (ms *MemoryStorage) GetAllMetrics() []models.Metric {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	return ms.metrics
}
