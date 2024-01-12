package logger

import (
	"go.uber.org/zap"
)

type logger2 struct {
	zap *zap.Logger
}

var Log2 = logger2{zap.NewNop()}

//
//// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
//func Initialize(level string) error {
//	// преобразуем текстовый уровень логирования в zap.AtomicLevel
//	lvl, err := zap.ParseAtomicLevel(level)
//	if err != nil {
//		return err
//	}
//	// создаём новую конфигурацию логера
//	cfg := zap.NewProductionConfig()
//
//	// устанавливаем уровень
//	cfg.Level = lvl
//	cfg.Encoding = "console"
//
//	l, err := cfg.Build()
//
//	if err != nil {
//		return err
//	}
//	// устанавливаем синглтон
//	Log = logger{l}
//
//	return nil
//}
//
//// RequestLogger — middleware-логер для входящих HTTP-запросов.
//func (l *logger) RequestLogger(h http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		log := Log.zap.With(
//			zap.String("remote_addr", r.RemoteAddr),
//		)
//		log.Info("started with: ", zap.String("method", r.Method), zap.String("uri", r.RequestURI))
//
//		responseData := &ResponseData{
//			status: 0,
//			size:   0,
//		}
//
//		lw := loggingResponseWriter{
//			ResponseWriter: w,
//			responseData:   responseData,
//		}
//		h.ServeHTTP(&lw, r)
//
//		start := time.Now()
//
//		var level zapcore.Level
//		switch {
//		case lw.responseData.status >= 500:
//			level = zap.ErrorLevel
//		case lw.responseData.status >= 400:
//			level = zap.WarnLevel
//		default:
//			level = zap.InfoLevel
//		}
//
//		log.Log(
//			level,
//			"completed with: ",
//			zap.Int("code", lw.responseData.status),
//			zap.Int("size", lw.responseData.size),
//			zap.String("status", http.StatusText(lw.responseData.status)),
//			zap.Duration("time", time.Since(start)),
//		)
//	})
//}
//
//func (l *logger) adapter(args ...interface{}) {
//}
//func (l *logger) Trace(args ...interface{}) {
//}
//func (l *logger) Debug(args ...interface{}) {
//}
//func (l *logger) Info(args ...interface{}) {
//}
//func (l *logger) Error(args ...interface{}) {
//}
//func (l *logger) Fatal(args ...interface{}) {
//}
