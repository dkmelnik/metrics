package logger

import (
	"context"
	"net/http"
)

type ILogger interface {
	RequestLog(h http.Handler) http.Handler
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	ErrorWithContext(ctx context.Context, err error)
}

type IConfig interface {
	GetLevel() string
	GetMode() string
}
