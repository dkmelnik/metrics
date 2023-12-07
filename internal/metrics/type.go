package metrics

type Kind string

// Возможные роли в семье.
const (
	Gauge   Kind = "gauge"
	Counter Kind = "counter"
)
