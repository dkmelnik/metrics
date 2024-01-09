package storage

import (
	"encoding/json"
	"errors"
	"github.com/dkmelnik/metrics/internal/models"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type MemoryStorage struct {
	mu          sync.RWMutex
	metrics     []models.Metric
	syncsSaving bool
	filePath    string
}

func NewMemoryStorage(storagePath string, storeInterval int, restore bool) (*MemoryStorage, error) {
	ms := &MemoryStorage{
		filePath:    storagePath,
		syncsSaving: storeInterval == 0,
		metrics:     make([]models.Metric, 0),
	}

	if restore {
		ms.loadMetricsFromFile()
	}

	if storeInterval > 0 {
		savePeriod := time.NewTicker(time.Second * time.Duration(storeInterval))
		go ms.intervalUpdatingToFile(savePeriod)
	}

	return ms, nil
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
	if ms.syncsSaving {
		ms.saveMetricsToFile()
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

func (ms *MemoryStorage) intervalUpdatingToFile(t *time.Ticker) {
	for {
		select {
		case <-t.C:
			ms.saveMetricsToFile()
		}
	}
}

func (ms *MemoryStorage) loadMetricsFromFile() {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	file, err := os.OpenFile(ms.filePath, os.O_RDONLY, 0666)
	if err != nil {
		log.Println("Error open file:", err)
		return
	}
	defer file.Close()

	ms.metrics = make([]models.Metric, 0)

	var m []models.Metric
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&m)
	if err != nil {
		log.Println("Error decode from file:", err)
		return
	}

	ms.metrics = m
}

func (ms *MemoryStorage) saveMetricsToFile() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if len(ms.metrics) > 0 {
		file, err := os.OpenFile(ms.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Println("Error opening file:", err)
			return
		}
		defer file.Close()
		encoder := json.NewEncoder(file)
		if err := encoder.Encode(ms.metrics); err != nil {
			log.Println("Error encoding metric:", err)
			return
		}
	}

}
