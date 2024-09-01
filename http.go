package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Метрики запросов
var (
	// HTTPRequestDuration отслеживает время, затраченное на обработку HTTP-запросов (в секундах).
	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	// HTTPRequestCount отслеживает количество обработанных HTTP-запросов.
	HTTPRequestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_count_total",
		Help: "Total number of HTTP requests processed",
	}, []string{"method", "status"})

	// HTTPErrorRate отслеживает процент ошибок среди всех запросов (4xx и 5xx ответы).
	HTTPErrorRate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_error_rate",
		Help: "Rate of HTTP request errors (4xx and 5xx responses)",
	}, []string{"method", "status"})

	// ResponseSize отслеживает размер ответов на HTTP-запросы (в байтах).
	ResponseSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_response_size_bytes",
		Help:    "Size of HTTP responses in bytes",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})
)

// Метрики обработки ошибок
var (
	// ErrorCount отслеживает количество ошибок в приложении.
	ErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "app_error_count_total",
		Help: "Total number of errors in the application",
	}, []string{"error_type"})

	// ErrorRate отслеживает процентное соотношение ошибок по отношению к общему количеству операций.
	ErrorRate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "app_error_rate",
		Help: "Percentage of errors relative to total operations",
	}, []string{"operation"})
)

// Метрики задержек
var (
	// Latency отслеживает задержки в выполнении критических операций (в секундах).
	Latency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "app_latency_seconds",
		Help:    "Latency of critical operations in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"operation"})
)

// Метрики зависимости
var (
	// ExternalServiceCallDuration отслеживает время выполнения вызовов к внешним сервисам (в секундах).
	ExternalServiceCallDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "external_service_call_duration_seconds",
		Help:    "Duration of calls to external services in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"service_name", "method"})

	// ExternalServiceErrorRate отслеживает процент ошибок при вызовах к внешним сервисам.
	ExternalServiceErrorRate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "external_service_error_rate",
		Help: "Error rate of calls to external services",
	}, []string{"service_name"})
)

// MeasureHTTPRequest Middleware для отслеживания метрик HTTP-запросов
func MeasureHTTPRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := &responseRecorder{w, http.StatusOK, 0}
		next.ServeHTTP(rr, r)

		duration := time.Since(start).Seconds()
		HTTPRequestDuration.WithLabelValues(r.Method, http.StatusText(rr.statusCode)).Observe(duration)
		HTTPRequestCount.WithLabelValues(r.Method, http.StatusText(rr.statusCode)).Inc()

		// Если статус код в диапазоне 4xx или 5xx, увеличиваем процент ошибок
		if rr.statusCode >= 400 && rr.statusCode < 600 {
			HTTPErrorRate.WithLabelValues(r.Method, http.StatusText(rr.statusCode)).Set(1)
		} else {
			HTTPErrorRate.WithLabelValues(r.Method, http.StatusText(rr.statusCode)).Set(0)
		}

		// Отслеживаем размер ответа
		ResponseSize.WithLabelValues(r.Method, http.StatusText(rr.statusCode)).Observe(float64(rr.responseSize))
	})
}

// responseRecorder используется для записи ответа и захвата статуса и размера ответа
type responseRecorder struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
}

// WriteHeader Переопределяем метод WriteHeader для записи статуса ответа
func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

// Write Переопределяем метод Write для записи размера ответа
func (rr *responseRecorder) Write(b []byte) (int, error) {
	size, err := rr.ResponseWriter.Write(b)
	rr.responseSize += size
	return size, err
}

// UpdateHTTPRequestDuration обновляет длительность обработки HTTP-запроса
func UpdateHTTPRequestDuration(method, status string, duration float64) {
	HTTPRequestDuration.WithLabelValues(method, status).Observe(duration)
}

// IncrementHTTPRequestCount увеличивает количество обработанных HTTP-запросов
func IncrementHTTPRequestCount(method, status string) {
	HTTPRequestCount.WithLabelValues(method, status).Inc()
}

// UpdateHTTPErrorRate обновляет процент ошибок для HTTP-запросов
func UpdateHTTPErrorRate(method, status string, rate float64) {
	HTTPErrorRate.WithLabelValues(method, status).Set(rate)
}

// ObserveResponseSize измеряет и обновляет размер ответа на HTTP-запрос
func ObserveResponseSize(method, status string, size float64) {
	ResponseSize.WithLabelValues(method, status).Observe(size)
}

// IncrementErrorCount увеличивает количество ошибок в приложении
func IncrementErrorCount(errorType string) {
	ErrorCount.WithLabelValues(errorType).Inc()
}

// UpdateErrorRate обновляет процент ошибок по отношению к общему количеству операций
func UpdateErrorRate(operation string, rate float64) {
	ErrorRate.WithLabelValues(operation).Set(rate)
}

// ObserveLatency измеряет и обновляет задержку выполнения критических операций
func ObserveLatency(operation string, duration float64) {
	Latency.WithLabelValues(operation).Observe(duration)
}

// UpdateExternalServiceCallDuration обновляет длительность вызова внешнего сервиса
func UpdateExternalServiceCallDuration(serviceName, method string, duration float64) {
	ExternalServiceCallDuration.WithLabelValues(serviceName, method).Observe(duration)
}

// UpdateExternalServiceErrorRate обновляет процент ошибок при вызовах внешних сервисов
func UpdateExternalServiceErrorRate(serviceName string, rate float64) {
	ExternalServiceErrorRate.WithLabelValues(serviceName).Set(rate)
}
