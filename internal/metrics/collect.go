package metrics

import (
	"github.com/dkmelnik/metrics/internal/models"
	"github.com/dkmelnik/metrics/internal/utils"
	"runtime"
	"time"
)

func Collect(md *models.Metrics) {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	pollCount := 0

	for {
		*md = models.Metrics{
			Alloc:         m.Alloc,
			TotalAlloc:    m.TotalAlloc,
			Sys:           m.Sys,
			Lookups:       m.Lookups,
			Mallocs:       m.Mallocs,
			Frees:         m.Frees,
			HeapAlloc:     m.HeapAlloc,
			HeapSys:       m.HeapSys,
			HeapIdle:      m.HeapIdle,
			HeapInuse:     m.HeapInuse,
			HeapReleased:  m.HeapReleased,
			HeapObjects:   m.HeapObjects,
			StackInuse:    m.StackInuse,
			StackSys:      m.StackSys,
			MSpanInuse:    m.MSpanInuse,
			MSpanSys:      m.MSpanSys,
			MCacheInuse:   m.MCacheInuse,
			MCacheSys:     m.MCacheSys,
			BuckHashSys:   m.BuckHashSys,
			GCSys:         m.GCSys,
			OtherSys:      m.OtherSys,
			NextGC:        m.NextGC,
			LastGC:        m.LastGC,
			PauseTotalNs:  m.PauseTotalNs,
			NumGC:         m.NumGC,
			NumForcedGC:   m.NumForcedGC,
			GCCPUFraction: m.GCCPUFraction,
			EnableGC:      m.EnableGC,
			DebugGC:       m.DebugGC,
			BySize:        m.BySize,
			PollCount:     pollCount,
			RandomValue:   utils.GenerateGuid(),
		}

		pollCount++
		time.Sleep(time.Second * 2)
	}
}
