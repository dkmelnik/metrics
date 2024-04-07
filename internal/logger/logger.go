package logger

import (
	"context"
	"net/http"
)

// ILogger is an interface defining logging functionalities.
type ILogger interface {
	RequestLog(h http.Handler) http.Handler
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	ErrorWithContext(ctx context.Context, err error)
}

// IConfig is an interface defining methods for retrieving configuration settings for init logger service.
type IConfig interface {
	GetLevel() string
	GetMode() string
}
