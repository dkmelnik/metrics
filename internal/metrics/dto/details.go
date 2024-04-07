package dto

import "github.com/dkmelnik/metrics/internal/models"

// Details represents details of an entity, including ID, type, delta, and value.
//
// The ID and MType fields are mandatory, while the Delta and Value fields are optional.
type Details struct {
	// ID is the unique identifier of the metric.
	ID string `json:"id"`
	// MType indicates the type of the metric, which can be either "gauge" or "counter".
	MType string `json:"type"`
	// Delta represents the change in value of the metric in case of a "counter" type.
	Delta *int64 `json:"delta,omitempty"`
	// Value represents the value of the metric in case of a "gauge" type.
	Value *float64 `json:"value,omitempty"`
}

// FillFromModel fills the Details struct with data from a given models.Metric model.
func (r *Details) FillFromModel(m models.Metric) {
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
