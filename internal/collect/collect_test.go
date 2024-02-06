package collect

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCollect(t *testing.T) {
	mockTicker := time.NewTicker(100 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())

	inputCh := MetricsGenerator(ctx, mockTicker)

	time.AfterFunc(1*time.Second, func() {
		mockTicker.Stop()
		cancel()
	})

	metricsReceived := make([]*Metrics, 0)
	for metric := range inputCh {
		metricsReceived = append(metricsReceived, metric)
	}

	assert.NotEmpty(t, metricsReceived, "Metrics need to be collected")
	assert.Equal(t, 10, len(metricsReceived), "At least 10 dimensions")
	assert.Equal(t, 10, metricsReceived[len(metricsReceived)-1].PollCount, "The number of collect updates should be equivalent")
}
