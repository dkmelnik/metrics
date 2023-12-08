package models

import "reflect"

type Metrics struct {
	Alloc         uint64  `metric:"gauge"`
	TotalAlloc    uint64  `metric:"gauge"`
	Sys           uint64  `metric:"gauge"`
	Lookups       uint64  `metric:"gauge"`
	Mallocs       uint64  `metric:"gauge"`
	Frees         uint64  `metric:"gauge"`
	HeapAlloc     uint64  `metric:"gauge"`
	HeapSys       uint64  `metric:"gauge"`
	HeapIdle      uint64  `metric:"gauge"`
	HeapInuse     uint64  `metric:"gauge"`
	HeapReleased  uint64  `metric:"gauge"`
	HeapObjects   uint64  `metric:"gauge"`
	StackInuse    uint64  `metric:"gauge"`
	StackSys      uint64  `metric:"gauge"`
	MSpanInuse    uint64  `metric:"gauge"`
	MSpanSys      uint64  `metric:"gauge"`
	MCacheInuse   uint64  `metric:"gauge"`
	MCacheSys     uint64  `metric:"gauge"`
	BuckHashSys   uint64  `metric:"gauge"`
	GCSys         uint64  `metric:"gauge"`
	OtherSys      uint64  `metric:"gauge"`
	NextGC        uint64  `metric:"gauge"`
	LastGC        uint64  `metric:"gauge"`
	PauseTotalNs  uint64  `metric:"gauge"`
	NumGC         uint32  `metric:"gauge"`
	NumForcedGC   uint32  `metric:"gauge"`
	GCCPUFraction float64 `metric:"gauge"`
	PollCount     int     `metric:"counter"`
	RandomValue   string  `metric:"gauge"`
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
