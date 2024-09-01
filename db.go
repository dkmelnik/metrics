package metrics

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Метрики соединений с базой данных
var (
	// OpenConnections отслеживает количество открытых соединений с базой данных.
	OpenConnections = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "db_open_connections",
		Help: "Number of open connections to the database",
	}, []string{"dbname"})

	// InUseConnections отслеживает количество соединений, которые в данный момент используются.
	InUseConnections = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "db_in_use_connections",
		Help: "Number of connections currently in use",
	}, []string{"dbname"})

	// IdleConnections отслеживает количество простаивающих соединений.
	IdleConnections = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "db_idle_connections",
		Help: "Number of idle connections",
	}, []string{"dbname"})

	// MaxOpenConnections отслеживает максимальное количество открытых соединений в пуле базы данных.
	MaxOpenConnections = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "db_max_open_connections",
		Help: "Maximum number of open connections in the database pool",
	}, []string{"dbname"})
)

// Метрики производительности запросов
var (
	// QueryDurationQuantile отслеживает длительность выполнения запросов, агрегированную в квантили.
	QueryDurationQuantile = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "db_query_duration_quantile_seconds",
		Help:       "Duration of database queries aggregated as quantiles",
		Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.005, 0.99: 0.001},
	}, []string{"query_name"})

	// QueryCount отслеживает общее количество выполненных запросов.
	QueryCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "db_query_count_total",
		Help: "Total number of database queries executed",
	}, []string{"query_name"})

	// SlowQueries отслеживает количество медленных запросов (превышающих заданный порог времени).
	SlowQueries = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "db_slow_queries_total",
		Help: "Number of slow database queries (exceeding threshold)",
	}, []string{"query_name"})
)

// Метрики ошибок
var (
	// DatabaseErrors отслеживает общее количество ошибок, связанных с базой данных.
	DatabaseErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "db_errors_total",
		Help: "Total number of database errors",
	}, []string{"db_name"})

	// ConnectionErrors отслеживает количество ошибок при установлении соединений с базой данных.
	ConnectionErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "db_connection_errors_total",
		Help: "Total number of errors when establishing connections to the database",
	}, []string{"db_name"})
)

// UpdateOpenConnections обновляет количество открытых соединений
func UpdateOpenConnections(dbname string, count float64) {
	OpenConnections.WithLabelValues(dbname).Set(count)
}

// UpdateInUseConnections обновляет количество используемых соединений
func UpdateInUseConnections(dbname string, count float64) {
	InUseConnections.WithLabelValues(dbname).Set(count)
}

// UpdateIdleConnections обновляет количество простаивающих соединений
func UpdateIdleConnections(dbname string, count float64) {
	IdleConnections.WithLabelValues(dbname).Set(count)
}

// UpdateMaxOpenConnections обновляет максимальное количество открытых соединений
func UpdateMaxOpenConnections(dbname string, count float64) {
	MaxOpenConnections.WithLabelValues(dbname).Set(count)
}

// ObserveQueryDuration измеряет и обновляет длительность выполнения запроса
func ObserveQueryDuration(start time.Time, queryName string, threshold ...float64) {
	duration := time.Since(start)

	QueryDurationQuantile.WithLabelValues(queryName).Observe(duration.Seconds())
	QueryCount.WithLabelValues(queryName).Inc()

	if len(threshold) == 0 {
		return
	}

	if duration.Seconds() > threshold[0] {
		IncrementSlowQueries(queryName)
	}
}

// IncrementSlowQueries увеличивает счетчик медленных запросов
func IncrementSlowQueries(queryName string) {
	SlowQueries.WithLabelValues(queryName).Inc()
}

// IncrementDatabaseErrors увеличивает счетчик ошибок базы данных
func IncrementDatabaseErrors(dbname string) {
	DatabaseErrors.WithLabelValues(dbname).Inc()
}

// IncrementConnectionErrors увеличивает счетчик ошибок подключения к базе данных
func IncrementConnectionErrors(dbname string) {
	ConnectionErrors.WithLabelValues(dbname).Inc()
}

// DBConnStat собирает метрики соединений с базой данных
func DBConnStat(ctx context.Context, t time.Duration, pool *pgxpool.Pool, dbname string) {
	go func() {
		ticker := time.NewTicker(t)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				stats := pool.Stat()

				UpdateOpenConnections(dbname, float64(stats.TotalConns()))
				UpdateInUseConnections(dbname, float64(stats.AcquiredConns()))
				UpdateIdleConnections(dbname, float64(stats.IdleConns()))
				UpdateMaxOpenConnections(dbname, float64(stats.MaxConns()))

			case <-ctx.Done():
				return
			}
		}
	}()
}
