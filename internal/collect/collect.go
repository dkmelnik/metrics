package collect

import (
	"context"
	"github.com/dkmelnik/metrics/internal/logger"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"math/rand"
	"runtime"
	"time"
)

func MetricsGenerator(ctx context.Context, t *time.Ticker) chan *Metrics {
	inputCh := make(chan *Metrics)
	var m runtime.MemStats
	pollCount := 0

	go func() {
		defer close(inputCh)

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				pollCount++

				memr, err := mem.VirtualMemory()
				if err != nil {
					logger.Log.Error("failed to get memory info", "error", err)
				}
				runtime.ReadMemStats(&m)

				cpuUsage, _ := cpu.Percent(0, false)
				numCPU := runtime.NumCPU()
				totalCPUUsage := 0.0
				for _, usage := range cpuUsage {
					totalCPUUsage += usage
				}
				avgCPUUsage := totalCPUUsage / float64(numCPU)

				inputCh <- &Metrics{
					Alloc:           float64(m.Alloc),
					TotalAlloc:      float64(m.TotalAlloc),
					Sys:             float64(m.Sys),
					Lookups:         float64(m.Lookups),
					Mallocs:         float64(m.Mallocs),
					Frees:           float64(m.Frees),
					HeapAlloc:       float64(m.HeapAlloc),
					HeapSys:         float64(m.HeapSys),
					HeapIdle:        float64(m.HeapIdle),
					HeapInuse:       float64(m.HeapInuse),
					HeapReleased:    float64(m.HeapReleased),
					HeapObjects:     float64(m.HeapObjects),
					StackInuse:      float64(m.StackInuse),
					StackSys:        float64(m.StackSys),
					MSpanInuse:      float64(m.MSpanInuse),
					MSpanSys:        float64(m.MSpanSys),
					MCacheInuse:     float64(m.MCacheInuse),
					MCacheSys:       float64(m.MCacheSys),
					BuckHashSys:     float64(m.BuckHashSys),
					GCSys:           float64(m.GCSys),
					OtherSys:        float64(m.OtherSys),
					NextGC:          float64(m.NextGC),
					LastGC:          float64(m.LastGC),
					PauseTotalNs:    float64(m.PauseTotalNs),
					NumGC:           float64(m.NumGC),
					NumForcedGC:     float64(m.NumForcedGC),
					TotalMemory:     float64(memr.Total),
					FreeMemory:      float64(memr.Free),
					CPUutilization1: avgCPUUsage,
					GCCPUFraction:   m.GCCPUFraction,
					PollCount:       pollCount,
					RandomValue:     rand.Float64(),
				}
			}
		}
	}()

	return inputCh
}
