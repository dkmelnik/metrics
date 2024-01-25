package storage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dkmelnik/metrics/internal/apperrors"
	"github.com/dkmelnik/metrics/internal/logger"
	"github.com/dkmelnik/metrics/internal/models"
	"github.com/dkmelnik/metrics/internal/utils"
	"os"
	"sync"
	"time"
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

	if restore {
		ms.loadMetricsFromFile()
	}

	if storeInterval > 0 {
		savePeriod := time.NewTicker(time.Second * time.Duration(storeInterval))
		go ms.intervalUpdatingToFile(savePeriod)
	}

	return ms, nil
}

func (m *MemoryStorage) SaveOrUpdate(ctx context.Context, metric models.Metric) error {
	existingMetric, err := m.FindOneByTypeAndName(ctx, metric.MType, metric.Name)

	m.mu.RLock()
	defer m.mu.RUnlock()

	if nil != err {
		if errors.Is(err, apperrors.ErrNotFound) {
			metric.ID = utils.GenerateGUID()
			m.metrics[metric.ID] = metric
			return nil
		}
		return err
	}

	metric.ID = existingMetric.ID
	m.metrics[existingMetric.ID] = metric

	return nil
}

func (m *MemoryStorage) SaveOrUpdateMany(ctx context.Context, metrics []models.Metric) error {
	for _, metric := range metrics {
		err := m.SaveOrUpdate(ctx, metric)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MemoryStorage) FindOneByTypeAndName(ctx context.Context, mType, mName string) (models.Metric, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, metric := range m.metrics {
		if metric.MType == mType && metric.Name == mName {
			return metric, nil
		}
	}

	return models.Metric{}, apperrors.ErrNotFound
}

func (m *MemoryStorage) Find(ctx context.Context) ([]models.Metric, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var existingData []models.Metric
	for _, metric := range m.metrics {
		existingData = append(existingData, metric)
	}

	return existingData, nil
}

func (m *MemoryStorage) intervalUpdatingToFile(t *time.Ticker) {
	for range t.C {
		m.saveMetricsToFile()
	}
}

func (m *MemoryStorage) loadMetricsFromFile() {
	ctx := context.Background()

	m.mu.Lock()
	defer m.mu.Unlock()

	var file *os.File
	var err error

	operation := func() error {
		file, err = os.OpenFile(m.filePath, os.O_RDONLY, 0666)
		return err
	}
	if err := apperrors.RetryWithBackoff(ctx, operation); err != nil {
		return
	}

	defer file.Close()

	var ms = make(map[string]models.Metric)

	err = json.NewDecoder(file).Decode(&ms)
	if err != nil {
		logger.Log.ErrorWithContext(ctx, err)
		return
	}
	for id, metric := range ms {
		key := id
		metric.ID = key
		m.metrics[key] = metric
	}
}

func (m *MemoryStorage) saveMetricsToFile() {
	ctx := context.Background()

	m.mu.Lock()
	defer m.mu.Unlock()

	var file *os.File
	var err error
	operation := func() error {
		file, err = os.OpenFile(m.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		return err
	}
	if err := apperrors.RetryWithBackoff(ctx, operation); err != nil {
		return
	}

	defer file.Close()
	encoder := json.NewEncoder(file)
	if err = encoder.Encode(m.metrics); err != nil {
		logger.Log.ErrorWithContext(ctx, err)
		return
	}
}
