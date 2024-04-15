package storage

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/dkmelnik/metrics/internal/models"
)

func BenchmarkSaveOrUpdate(b *testing.B) {
	storage, _ := NewMemoryStorage("", 0, false)

	metric := models.Metric{
		MType: "type",
		Name:  "name",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = storage.SaveOrUpdate(context.Background(), metric)
	}
}

func BenchmarkSaveOrUpdateMany(b *testing.B) {
	storage, _ := NewMemoryStorage("", 0, false)

	var metrics []models.Metric
	for i := 0; i < 1000; i++ {
		metric := models.Metric{
			MType: "type",
			Name:  "name",
		}
		metrics = append(metrics, metric)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = storage.SaveOrUpdateMany(context.Background(), metrics)
	}
}

func BenchmarkFindOneByTypeAndName(b *testing.B) {
	storage, _ := NewMemoryStorage("", 0, false)

	for i := 0; i < 1000; i++ {
		metric := models.Metric{
			MType: "type",
			Name:  "name" + strconv.Itoa(i),
		}
		storage.metrics[metric.ID] = metric
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = storage.FindOneByTypeAndName(context.Background(), "type", "name"+strconv.Itoa(i))
	}
}

func BenchmarkFind(b *testing.B) {
	storage, _ := NewMemoryStorage("", 0, false)

	for i := 0; i < 1000; i++ {
		metric := models.Metric{
			MType: "type",
			Name:  "name" + strconv.Itoa(i),
		}
		storage.metrics[metric.ID] = metric
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = storage.Find(context.Background())
	}
}

// -------------------------------------------
func generateSampleMetricsFile(filePath string) error {
	metrics := make(map[string]models.Metric)
	for i := 0; i < 1000; i++ {
		metric := models.Metric{
			ID:    strconv.Itoa(i),
			MType: "type",
			Name:  "name" + strconv.Itoa(i),
		}
		metrics[strconv.Itoa(i)] = metric
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(metrics); err != nil {
		return err
	}

	return nil
}

func BenchmarkLoadMetricsFromFile(b *testing.B) {
	tmpFilePath := "tmp_metrics.json"
	defer os.Remove(tmpFilePath)

	if err := generateSampleMetricsFile(tmpFilePath); err != nil {
		b.Fatal("Error generating sample metrics file:", err)
	}

	storage, _ := NewMemoryStorage(tmpFilePath, 0, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.loadMetricsFromFile()
	}
}

// -------------------------------------------
func BenchmarkSaveMetricsToFile(b *testing.B) {
	// Create a temporary file for benchmarking
	tmpFilePath := "tmp_metrics.json"
	defer os.Remove(tmpFilePath)

	// Initialize MemoryStorage
	storage, _ := NewMemoryStorage(tmpFilePath, 0, false)

	// Populate storage with some sample metrics
	for i := 0; i < 1000; i++ {
		metric := models.Metric{
			ID:    strconv.Itoa(i),
			MType: "type",
			Name:  "name" + strconv.Itoa(i),
			// Set other fields as needed for your benchmark
		}
		storage.metrics[strconv.Itoa(i)] = metric
	}

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Call the method to be benchmarked
		storage.saveMetricsToFile()
	}
}

// -------------------------------------------
