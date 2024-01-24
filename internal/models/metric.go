package models

import (
	"database/sql"
	"github.com/dkmelnik/metrics/internal/apperrors"
	"math"
	"time"
)

type MetricType string

var (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type Metric struct {
	ID        string          `json:"id" db:"id"`                 // имя метрики
	Name      string          `json:"name" db:"name"`             // параметр, принимающий значение gauge или counter
	MType     string          `json:"type" db:"type"`             // параметр, принимающий значение gauge или counter
	Delta     sql.NullInt64   `json:"delta,omitempty" db:"delta"` // значение метрики в случае передачи counter
	Value     sql.NullFloat64 `json:"value,omitempty" db:"value"` // значение метрики в случае передачи gauge
	CreatedAT time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAT time.Time       `json:"updatedAt" db:"updated_at"`
}

func NewMetric(name string, mType string) (Metric, error) {
	if !(mType == string(Gauge) || mType == string(Counter)) {
		return Metric{}, apperrors.ErrTypeNotCorrect
	}

	return Metric{
		Name:  name,
		MType: mType,
	}, nil
}

func (m *Metric) SetDelta(delta int64) {
	m.Delta.Int64 = delta
	m.Delta.Valid = true
}

func (m *Metric) UpdateDelta(delta int64) {
	m.Delta.Int64 += delta
}

func (m *Metric) SetValue(value float64) {
	m.Value.Float64 = value
	m.Value.Valid = true
}

func (m *Metric) GetAdaptedByTypeValue() interface{} {
	switch m.MType {
	case string(Counter):
		return m.Delta.Int64
	default:
		return m.Value.Float64
	}
}

func (m *Metric) round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func (m *Metric) toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(m.round(num*output)) / output
}

func (m *Metric) CheckType() error {
	if !(m.MType == string(Gauge) || m.MType == string(Counter)) {
		return apperrors.ErrTypeNotCorrect
	}
	return nil
}
