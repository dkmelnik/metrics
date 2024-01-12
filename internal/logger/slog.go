package logger

import (
	"context"
	"encoding/json"
	"github.com/dkmelnik/metrics/configs"
	"github.com/fatih/color"
	"github.com/mdobak/go-xerrors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"
)

/*

type User struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	logger.Log.Debug("debug", "user", "user", "user")
	logger.Log.Info("user", "user", "user")
	logger.Log.Warn("warn", "user", "user")
	logger.Log.Error("error", "user", "user")
	u := &User{
		ID:        "user-12234",
		FirstName: "Jan",
		LastName:  "Doe",
		Email:     "jan@example.com",
		Password:  "pass-12334",
	}
	logger.Log.Error("user", "user", u)
	logger.Log.ErrorWithContext(context.TODO(), errors.New("some error"))
*/

type logger struct {
	slog *slog.Logger
}

var _ ILogger = (*logger)(nil)

var Log = logger{slog.New(slog.NewJSONHandler(io.Discard, nil))}

func (l *logger) Debug(msg string, args ...any) {
	l.slog.Debug(msg, args...)
}
func (l *logger) Info(msg string, args ...any) {
	l.slog.Info(msg, args...)
}
func (l *logger) Warn(msg string, args ...any) {
	l.slog.Warn(msg, args...)
}
func (l *logger) Error(msg string, args ...any) {
	l.slog.Error(msg, args...)
}
func (l *logger) RequestLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.slog.Info("incoming request: ", slog.String("method", r.Method), slog.String("uri", r.RequestURI))

		responseData := &ResponseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		start := time.Now()

		var level slog.Level
		switch {
		case lw.responseData.status >= 500:
			level = slog.LevelError
		case lw.responseData.status >= 400:
			level = slog.LevelWarn
		default:
			level = slog.LevelInfo
		}
		l.slog.Log(
			r.Context(),
			level,
			"completed with: ",
			slog.Int("code", lw.responseData.status),
			slog.Int("size", lw.responseData.size),
			slog.String("status", http.StatusText(lw.responseData.status)),
			slog.Duration("time", time.Since(start)),
		)
	})
}
func (l *logger) ErrorWithContext(ctx context.Context, err error) {
	var sl slog.Attr
	if err != nil {
		err = xerrors.New(err.Error())
		sl = slog.Any("error", err)
	}
	l.slog.ErrorContext(ctx, err.Error(), sl)
}

func Setup(c configs.Logger, w io.Writer) error {
	handler := setupHandler(c, w)

	Log = logger{slog.New(handler)}
	return nil
}

// logger level mapping
var loggerLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func getLoggerLevel(c configs.Logger) slog.Level {
	level, exist := loggerLevelMap[c.Level]
	if !exist {
		return slog.LevelError
	}

	return level
}

// Depending on the operating mode of the application, select the logger operating mode
func setupHandler(c configs.Logger, w io.Writer) slog.Handler {
	opts := &slog.HandlerOptions{
		Level: getLoggerLevel(c),
	}
	if c.Mode == "production" {
		opts.ReplaceAttr = replaceAttr
		return slog.NewJSONHandler(w, opts)
	}

	return newPrettyHandler(w, opts)
}

// For dev mode of the application
type prettyHandler struct {
	slog.Handler
	l *log.Logger
}

func newPrettyHandler(
	out io.Writer,
	opts *slog.HandlerOptions,
) *prettyHandler {
	h := &prettyHandler{
		Handler: slog.NewJSONHandler(out, opts),
		l:       log.New(out, "", 0),
	}

	return h
}

func (h *prettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		if a.Value.Kind() == slog.KindAny {
			switch v := a.Value.Any().(type) {
			case error:
				fields[a.Key] = map[string]interface{}{
					"msg":   v.Error(),
					"trace": marshalStack(v),
				}
			default:
				fields[a.Key] = a.Value.Any()
			}
		} else {
			fields[a.Key] = a.Value.Any()
		}

		return true
	})
	b, err := json.MarshalIndent(fields, "", "  ")

	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(timeStr, level, msg, color.WhiteString(string(b)))

	return nil
}

type stackFrame struct {
	Func   string `json:"func"`
	Source string `json:"source"`
	Line   int    `json:"line"`
}

func replaceAttr(_ []string, a slog.Attr) slog.Attr {
	switch a.Value.Kind() {
	case slog.KindAny:
		switch v := a.Value.Any().(type) {
		case error:
			a.Value = fmtErr(v)
		}
	}

	return a
}

// marshalStack extracts stack frames from the error
func marshalStack(err error) []stackFrame {
	trace := xerrors.StackTrace(err)

	if len(trace) == 0 {
		return nil
	}

	frames := trace.Frames()

	s := make([]stackFrame, len(frames))

	for i, v := range frames {
		f := stackFrame{
			Source: filepath.Join(
				filepath.Base(filepath.Dir(v.File)),
				filepath.Base(v.File),
			),
			Func: filepath.Base(v.Function),
			Line: v.Line,
		}

		s[i] = f
	}

	return s
}

// fmtErr returns a slog.Value with keys `msg` and `trace`. If the error
// does not implement interface { StackTrace() errors.StackTrace }, the `trace`
// key is omitted.
func fmtErr(err error) slog.Value {
	var groupValues []slog.Attr

	groupValues = append(groupValues, slog.String("msg", err.Error()))

	frames := marshalStack(err)
	if frames != nil {
		groupValues = append(groupValues,
			slog.Any("trace", frames),
		)
	}

	return slog.GroupValue(groupValues...)
}
