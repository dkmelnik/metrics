package apperrors

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/dkmelnik/metrics/internal/logger"
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
	case isTestError(err):
		return true, initialRetryDelay
	default:
		return false, 0
	}
}

func isPostgreSQLError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
}

func isTestError(err error) bool {
	return strings.Contains(err.Error(), "test")
}

func RetryWithBackoff(ctx context.Context, operation func() error) error {
	for attempt := 1; ; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}

		ok, delay := IsRetriableError(err)
		logger.Log.Debug(
			"RetryWithBackoff",
			"attempt",
			attempt,
		)

		if !ok || attempt >= maxRetries {
			return err
		}

		backoff := time.Duration(attempt+2) * delay
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
