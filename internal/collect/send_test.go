package collect

import (
	"context"
	"encoding/json"
	"github.com/dkmelnik/metrics/internal/models"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"time"
)

func Test_Send(t *testing.T) {
	metricsNames := make([]string, 0)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data models.Metric

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Println("Error decoding request body:", err)
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}

		metricsNames = append(metricsNames, data.ID)
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	md := &Metrics{}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	*md = Metrics{
		Alloc:         float64(m.Alloc),
		TotalAlloc:    float64(m.TotalAlloc),
		Sys:           float64(m.Sys),
		Lookups:       float64(m.Lookups),
		Mallocs:       float64(m.Mallocs),
		Frees:         float64(m.Frees),
		HeapAlloc:     float64(m.HeapAlloc),
		HeapSys:       float64(m.HeapSys),
		HeapIdle:      float64(m.HeapIdle),
		HeapInuse:     float64(m.HeapInuse),
		HeapReleased:  float64(m.HeapReleased),
		HeapObjects:   float64(m.HeapObjects),
		StackInuse:    float64(m.StackInuse),
		StackSys:      float64(m.StackSys),
		MSpanInuse:    float64(m.MSpanInuse),
		MSpanSys:      float64(m.MSpanSys),
		MCacheInuse:   float64(m.MCacheInuse),
		MCacheSys:     float64(m.MCacheSys),
		BuckHashSys:   float64(m.BuckHashSys),
		GCSys:         float64(m.GCSys),
		OtherSys:      float64(m.OtherSys),
		NextGC:        float64(m.NextGC),
		LastGC:        float64(m.LastGC),
		PauseTotalNs:  float64(m.PauseTotalNs),
		NumGC:         float64(m.NumGC),
		NumForcedGC:   float64(m.NumForcedGC),
		GCCPUFraction: m.GCCPUFraction,
		PollCount:     1,
		RandomValue:   rand.Float64(),
	}

	sendPeriod := time.NewTicker(time.Millisecond * 50)
	defer sendPeriod.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go Send(ctx, sendPeriod, md, server.URL)

	time.Sleep(time.Millisecond * 98)

	assert.ElementsMatch(t, md.GetProperties(), metricsNames, "each of the collect must be sent")
}

func Test_BuildRequestBody(t *testing.T) {
	testCases := []struct {
		desc     string
		metric   string
		name     string
		value    interface{}
		expected map[string]interface{}
	}{
		{
			desc:   "Build gauge request body",
			metric: "gauge",
			name:   "some_metric_name",
			value:  25.5,
			expected: map[string]interface{}{
				"id":    "some_metric_name",
				"type":  "gauge",
				"value": 25.5,
			},
		},
		{
			desc:   "Build counter request body",
			metric: "counter",
			name:   "another_metric_name",
			value:  10,
			expected: map[string]interface{}{
				"id":    "another_metric_name",
				"type":  "counter",
				"delta": 10,
			},
		},
		{
			desc:   "Invalid value type",
			metric: "gauge",
			name:   "invalid_metric",
			value:  "invalid_value",
			expected: map[string]interface{}{
				"id":   "invalid_metric",
				"type": "gauge",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := buildRequestBody(tc.metric, tc.name, tc.value)

			expectedJSON, _ := json.Marshal(tc.expected)
			resultJSON, _ := json.Marshal(result)

			assert.JSONEq(t, string(expectedJSON), string(resultJSON))
		})
	}
}
