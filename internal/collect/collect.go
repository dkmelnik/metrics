package collect

import (
	"context"
	"math/rand"
	"runtime"
	"time"
)

func Collect(ctx context.Context, t *time.Ticker, ch chan<- *Metrics) {
	var m runtime.MemStats
	pollCount := 0

	for {
		select {
		case <-ctx.Done():
			close(ch)
			return
		case <-t.C:
			pollCount++

			runtime.ReadMemStats(&m)
			ch <- &Metrics{
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
				PollCount:     pollCount,
				RandomValue:   rand.Float64(),
			}
		}
	}
}
