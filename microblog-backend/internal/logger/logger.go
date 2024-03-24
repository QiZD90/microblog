package logger

import (
	"context"
	"net/http"
	"os"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type ctxKey struct{}

type Logger interface {
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	WithFields(fields map[string]any) Logger
}

type logger struct {
	logrusLogger *logrus.Logger
	fields       map[string]any
}

func (log logger) Errorf(format string, args ...interface{}) {
	log.logrusLogger.WithFields(logrus.Fields(log.fields)).Errorf(format, args...)
}

func (log logger) Infof(format string, args ...interface{}) {
	log.logrusLogger.WithFields(logrus.Fields(log.fields)).Infof(format, args...)
}

func (log logger) WithFields(fields map[string]any) Logger {
	fieldsCopy := lo.MapEntries(log.fields, func(k string, v any) (string, any) { return k, v })
	for k, v := range fields {
		fieldsCopy[k] = v
	}

	return logger{
		logrusLogger: log.logrusLogger,
		fields:       fieldsCopy,
	}
}

func New() Logger {
	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})
	logrusLogger.SetOutput(os.Stdout)

	return logger{
		logrusLogger: logrusLogger,
		fields:       make(map[string]any),
	}
}

func NewContext(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, log)
}

func FromContext(ctx context.Context) Logger {
	switch log := ctx.Value(ctxKey{}).(type) {
	case Logger:
		return log
	default:
		return logStub{}
	}
}

func Inject(log Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(NewContext(r.Context(), log)))
		})
	}
}
