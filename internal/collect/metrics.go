package collect

import "reflect"

// Metrics represents system metrics including memory usage, CPU utilization, and other statistics.
//
// Each field in the Metrics struct corresponds to a specific metric, with its data type and metric type indicated using struct tags.
// Metric types can be "gauge" or "counter", indicating whether the metric represents a value that can go up and down (gauge) or only up (counter).
// The struct tags are used for metric annotations.
type Metrics struct {
	Alloc           float64 `metric:"gauge"`
	TotalAlloc      float64 `metric:"gauge"`
	Sys             float64 `metric:"gauge"`
	Lookups         float64 `metric:"gauge"`
	Mallocs         float64 `metric:"gauge"`
	Frees           float64 `metric:"gauge"`
	HeapAlloc       float64 `metric:"gauge"`
	HeapSys         float64 `metric:"gauge"`
	HeapIdle        float64 `metric:"gauge"`
	HeapInuse       float64 `metric:"gauge"`
	HeapReleased    float64 `metric:"gauge"`
	HeapObjects     float64 `metric:"gauge"`
	StackInuse      float64 `metric:"gauge"`
	StackSys        float64 `metric:"gauge"`
	MSpanInuse      float64 `metric:"gauge"`
	MSpanSys        float64 `metric:"gauge"`
	MCacheInuse     float64 `metric:"gauge"`
	MCacheSys       float64 `metric:"gauge"`
	BuckHashSys     float64 `metric:"gauge"`
	GCSys           float64 `metric:"gauge"`
	OtherSys        float64 `metric:"gauge"`
	NextGC          float64 `metric:"gauge"`
	LastGC          float64 `metric:"gauge"`
	PauseTotalNs    float64 `metric:"gauge"`
	NumGC           float64 `metric:"gauge"`
	NumForcedGC     float64 `metric:"gauge"`
	GCCPUFraction   float64 `metric:"gauge"`
	TotalMemory     float64 `metric:"gauge"`
	FreeMemory      float64 `metric:"gauge"`
	CPUutilization1 float64 `metric:"gauge"`
	PollCount       int     `metric:"counter"`
	RandomValue     float64 `metric:"gauge"`
}

func (m Metrics) GetProperties() []string {
	metricsType := reflect.TypeOf(m)

	var properties []string

	for i := 0; i < metricsType.NumField(); i++ {
		field := metricsType.Field(i)
		properties = append(properties, field.Name)
	}

	return properties
}
