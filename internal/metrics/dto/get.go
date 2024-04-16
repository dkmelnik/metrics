package dto

// GetRequest represents a request for retrieving an entity.
//
// It contains fields for specifying the ID and type of the entity to be retrieved.
type GetRequest struct {
	// ID is the unique identifier of the metric.
	ID string `json:"id"`
	// MType indicates the type of the metric, which can be either "gauge" or "counter".
	MType string `json:"type"`
}
