package collect

import (
	"context"
	"github.com/dkmelnik/metrics/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCollect(t *testing.T) {
	mockMetrics := &models.Metrics{}
	ctx, cancel := context.WithCancel(context.Background())

	go Collect(ctx, time.NewTicker(time.Second*2), mockMetrics)

	time.Sleep(time.Second * 5)
	cancel()

	assert.Equal(t, 2, mockMetrics.PollCount, "The number of collect updates should be equivalent")
}
