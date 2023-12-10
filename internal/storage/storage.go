package storage

import (
	"errors"
	"sync"
)

type MemoryStorage struct {
	mu      sync.RWMutex
	metrics map[string]map[string]interface{}
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		metrics: make(map[string]map[string]interface{}),
	}
}

func (ms *MemoryStorage) Save(metricType, metricName string, value interface{}) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if ms.metrics[metricType] == nil {
		ms.metrics[metricType] = make(map[string]interface{})
	}
	ms.metrics[metricType][metricName] = value
}

func (ms *MemoryStorage) FindOneByTypeName(metricType, metricName string) (interface{}, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if ms.metrics[metricType] == nil || ms.metrics[metricType][metricName] == 0 {
		return 0, errors.New("metric not found")
	}

	return ms.metrics[metricType][metricName], nil
}

func (ms *MemoryStorage) GetAllMetrics() map[string]map[string]interface{} {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	return ms.metrics
}
