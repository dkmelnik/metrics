package dto

import "github.com/dkmelnik/metrics/internal/models"

type GetRequest struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}

type Response struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (r *Response) AdaptModel(m models.Metric) {
	r.ID = m.Name
	r.MType = m.MType
	r.Delta = nil
	r.Value = nil

	if m.Delta.Valid {
		r.Delta = &m.Delta.Int64
	}

	if m.Value.Valid {
		r.Value = &m.Value.Float64
	}
}
