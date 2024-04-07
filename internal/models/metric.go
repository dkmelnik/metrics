package models

import (
	"database/sql"
	"math"
	"time"

	"github.com/dkmelnik/metrics/internal/apperrors"
)

// MetricType represents the type of a metric, which can be either "gauge" or "counter".
type MetricType string

// Constants representing the two possible types of metrics.
const (
	// Gauge represents a metric type that measures a value at a particular point in time.
	Gauge MetricType = "gauge"
	// Counter represents a metric type that measures a continuously increasing value over time.
	Counter MetricType = "counter"
)

// Metric represents a data structure for storing metric information.
type Metric struct {
	// ID is the unique identifier of the metric.
	ID string `json:"id" db:"id"`
	// Name is the name of the metric.
	Name string `json:"name" db:"name"`
	// MType indicates the type of the metric, which can be either "gauge" or "counter".
	MType string `json:"type" db:"type"`
	// Delta represents the change in value of the metric in case of a "counter" type.
	Delta sql.NullInt64 `json:"delta,omitempty" db:"delta"`
	// Value represents the value of the metric in case of a "gauge" type.
	Value sql.NullFloat64 `json:"value,omitempty" db:"value"`
	// CreatedAt represents the timestamp when the metric was created.
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	// UpdatedAt represents the timestamp when the metric was last updated.
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
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

func (m *Metric) GetValueByType() interface{} {
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
