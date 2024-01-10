package metrics

// TODO перенести в папку dto
type GetMetricRequest struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}
