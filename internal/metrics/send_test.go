package metrics

import (
	"context"
	"github.com/dkmelnik/metrics/internal/models"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	metricsNames := make([]string, 0)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		parts := strings.Split(r.URL.Path, "/")
		metricsNames = append(metricsNames, parts[3])
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	md := &models.Metrics{}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	*md = models.Metrics{
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

	assert.ElementsMatch(t, md.GetProperties(), metricsNames, "each of the metrics must be sent")
}

func TestBuildRequestURL(t *testing.T) {
	serverURL := "http://example.com"
	tag := "gauge"
	fieldName := "Alloc"
	value := "243288"

	expectedURL := "http://example.com/update/gauge/Alloc/243288"
	result := buildRequestURL(serverURL, tag, fieldName, value)

	assert.Equal(t, expectedURL, result)
}
