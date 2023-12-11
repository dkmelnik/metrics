package models

import "reflect"

type MetricType string

var (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type Metrics struct {
	Alloc         float64 `metric:"gauge"`
	TotalAlloc    float64 `metric:"gauge"`
	Sys           float64 `metric:"gauge"`
	Lookups       float64 `metric:"gauge"`
	Mallocs       float64 `metric:"gauge"`
	Frees         float64 `metric:"gauge"`
	HeapAlloc     float64 `metric:"gauge"`
	HeapSys       float64 `metric:"gauge"`
	HeapIdle      float64 `metric:"gauge"`
	HeapInuse     float64 `metric:"gauge"`
	HeapReleased  float64 `metric:"gauge"`
	HeapObjects   float64 `metric:"gauge"`
	StackInuse    float64 `metric:"gauge"`
	StackSys      float64 `metric:"gauge"`
	MSpanInuse    float64 `metric:"gauge"`
	MSpanSys      float64 `metric:"gauge"`
	MCacheInuse   float64 `metric:"gauge"`
	MCacheSys     float64 `metric:"gauge"`
	BuckHashSys   float64 `metric:"gauge"`
	GCSys         float64 `metric:"gauge"`
	OtherSys      float64 `metric:"gauge"`
	NextGC        float64 `metric:"gauge"`
	LastGC        float64 `metric:"gauge"`
	PauseTotalNs  float64 `metric:"gauge"`
	NumGC         float64 `metric:"gauge"`
	NumForcedGC   float64 `metric:"gauge"`
	GCCPUFraction float64 `metric:"gauge"`
	PollCount     int     `metric:"counter"`
	RandomValue   float64 `metric:"gauge"`
}

func (m Metrics) HasProperty(prop string) bool {
	structValue := reflect.ValueOf(m)

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Type().Field(i)
		if field.Name == prop {
			return true
		}
	}
	return false
}

func (m Metrics) HasType(tag string) bool {
	structType := reflect.TypeOf(m)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.Tag.Get("metric") == tag {
			return true
		}
	}
	return false
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
