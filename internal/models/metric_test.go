package models

import (
	"testing"
)

func TestNewMetric(t *testing.T) {
	name := "test_metric"
	mType := "gauge"

	metric, err := NewMetric(name, mType)
	if err != nil {
		t.Errorf("NewMetric() returned an error: %v", err)
	}

	if metric.Name != name {
		t.Errorf("NewMetric() returned a metric with incorrect name: got %s, want %s", metric.Name, name)
	}

	if metric.MType != mType {
		t.Errorf("NewMetric() returned a metric with incorrect type: got %s, want %s", metric.MType, mType)
	}
}

func TestMetric_SetDelta(t *testing.T) {
	metric := Metric{}
	delta := int64(10)

	metric.SetDelta(delta)

	if metric.Delta.Int64 != delta {
		t.Errorf("SetDelta() did not set the delta correctly: got %d, want %d", metric.Delta.Int64, delta)
	}

	if !metric.Delta.Valid {
		t.Error("SetDelta() did not set the delta as valid")
	}
}

func TestMetric_UpdateDelta(t *testing.T) {
	metric := Metric{}
	initialDelta := int64(5)
	delta := int64(10)

	metric.SetDelta(initialDelta)
	metric.UpdateDelta(delta)

	expectedDelta := initialDelta + delta
	if metric.Delta.Int64 != expectedDelta {
		t.Errorf("UpdateDelta() did not update the delta correctly: got %d, want %d", metric.Delta.Int64, expectedDelta)
	}
}

func TestMetric_SetValue(t *testing.T) {
	metric := Metric{}
	value := 3.14

	metric.SetValue(value)

	if metric.Value.Float64 != value {
		t.Errorf("SetValue() did not set the value correctly: got %f, want %f", metric.Value.Float64, value)
	}

	if !metric.Value.Valid {
		t.Error("SetValue() did not set the value as valid")
	}
}

func TestMetric_GetValueByType(t *testing.T) {
	metric := Metric{}
	metric.MType = "counter"
	delta := int64(10)
	metric.SetDelta(delta)

	value := metric.GetValueByType()
	if val, ok := value.(int64); ok {
		if val != delta {
			t.Errorf("GetValueByType() did not return the correct value for counter type: got %d, want %d", val, delta)
		}
	} else {
		t.Error("GetValueByType() did not return an integer value for counter type")
	}
}

func TestMetric_CheckType(t *testing.T) {
	metric := Metric{MType: "gauge"}

	err := metric.CheckType()
	if err != nil {
		t.Errorf("CheckType() returned an unexpected error: %v", err)
	}
}
