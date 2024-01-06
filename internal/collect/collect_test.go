package collect

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestCollect(t *testing.T) {
	metricsChan := make(chan *Metrics)

	mockTicker := time.NewTicker(100 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())

	go Collect(ctx, mockTicker, metricsChan)

	time.Sleep(800 * time.Millisecond)

	mockTicker.Stop()
	cancel()

	metricsReceived := make([]*Metrics, 0)
	for metric := range metricsChan {
		log.Printf("%v", metric)
		metricsReceived = append(metricsReceived, metric)
	}
	//assert.NotEmpty(t, metricsReceived, "Metrics need to be collected")
	//assert.Equal(t, 5, 2, "At least 3 dimensions")
	//assert.Equal(t, 5, metricsReceived[len(metricsReceived)].PollCount, "The number of collect updates should be equivalent")
}
