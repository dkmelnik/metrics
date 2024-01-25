package apperrors

import (
	"context"
	"errors"
	"github.com/dkmelnik/metrics/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

const (
	maxRetries        = 3
	initialRetryDelay = time.Second
	maxRetryDelay     = 5 * time.Second
)

func IsRetriableError(err error) (bool, time.Duration) {
	switch {
	case isPostgreSQLError(err):
		return true, initialRetryDelay
	default:
		return false, 0
	}
}

func isPostgreSQLError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
}

func RetryWithBackoff(ctx context.Context, operation func() error) error {
	for attempt := 1; ; attempt++ {
		err := operation()

		ok, delay := IsRetriableError(err)

		if !ok || attempt >= maxRetries {
			return err
		}

		logger.Log.Info("retryWithBackoff", "Retriable error occurred, retrying attempt %d\n", attempt)

		backoff := time.Duration(attempt+1) * delay
		if backoff > maxRetryDelay {
			backoff = maxRetryDelay
		}

		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
