package logger

import "net/http"

type Logger interface {
	RequestLogger(h http.Handler) http.Handler
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}
