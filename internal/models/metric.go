package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"time"
)

type MetricType string

var (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type Metric struct {
	ID        string          `json:"-" db:"id"`                  // имя метрики
	Name      string          `json:"id" db:"name"`               // параметр, принимающий значение gauge или counter
	MType     string          `json:"type" db:"type"`             // параметр, принимающий значение gauge или counter
	Delta     sql.NullInt64   `json:"delta,omitempty" db:"delta"` // значение метрики в случае передачи counter
	Value     sql.NullFloat64 `json:"value,omitempty" db:"value"` // значение метрики в случае передачи gauge
	CreatedAT time.Time       `json:"-" db:"created_at"`
	UpdatedAT time.Time       `json:"-" db:"updated_at"`
}

func (m *Metric) SetDelta(delta int64) {
	m.Delta.Int64 = delta
	m.Delta.Valid = true
}

func (m *Metric) UpdateDelta(delta int64) {
	m.Delta.Int64 += delta
}

func (m *Metric) SetValue(value float64) {
	m.Value.Float64 = m.toFixed(value, 3)
	m.Value.Valid = true
}

func (m *Metric) GetDelta() int64 {
	return m.Delta.Int64
}

func (m *Metric) GetValue() interface{} {
	switch m.MType {
	case string(Counter):
		return m.Delta.Int64
	default:
		return m.Value.Float64
	}
}

func (m Metric) round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func (m Metric) toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(m.round(num*output)) / output
}

func (m *Metric) UnmarshalJSON(data []byte) error {
	var temp struct {
		ID    string   `json:"id"`
		MType string   `json:"type"`
		Delta *int64   `json:"delta,omitempty"`
		Value *float64 `json:"value,omitempty"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	var d sql.NullInt64
	if temp.Delta != nil {
		d.Int64 = *temp.Delta
		d.Valid = true
	}
	var v sql.NullFloat64
	if temp.Value != nil {
		v.Float64 = m.toFixed(*temp.Value, 3)
		v.Valid = true
	}

	m.Name = temp.ID
	m.MType = temp.MType
	m.Delta = d
	m.Value = v

	return nil
}

func (m Metric) MarshalJSON() ([]byte, error) {
	t := struct {
		ID    string  `json:"id"`
		MType string  `json:"type"`
		Delta int64   `json:"delta,omitempty"`
		Value float64 `json:"value,omitempty"`
	}{
		ID:    m.Name,
		MType: m.MType,
	}
	switch m.MType {
	case string(Counter):
		t.Delta = m.Delta.Int64
	case string(Gauge):
		t.Value = m.Value.Float64
	default:
		return nil, fmt.Errorf("unsupported metric type: %s", m.MType)
	}

	out, err := json.Marshal(t)

	if err != nil {
		return nil, err
	}
	return out, nil
}
