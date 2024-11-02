package logger

import (
	"log"
	"log/slog"
	"os"
)

const (
	LevelDebug   = "debug"
	LevelInfo    = "info"
	LevelWarning = "warn"
	LevelError   = "error"
)

type Logger interface {
	WithField(key string, val string) Logger
	WithFields(fields map[string]any) Logger
	WithError(err error) Logger
	Error(msg string)
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
}

type ConfigOptions func(*slog.HandlerOptions)

func WithMinLogLevel(minErrorLevel string) ConfigOptions {
	return func(o *slog.HandlerOptions) {
		availableErrorLevels := map[string]struct{}{
			LevelDebug:   {},
			LevelInfo:    {},
			LevelWarning: {},
			LevelError:   {},
		}
		if _, ok := availableErrorLevels[minErrorLevel]; !ok {
			log.Fatalf("unsupported log level %s", minErrorLevel)
		}

		var level slog.Level
		err := level.UnmarshalText([]byte(minErrorLevel))
		if err != nil {
			log.Fatal(err)
		}

		o.Level = level
	}
}

type Log struct {
	original *slog.Logger
}

func New(opts ...ConfigOptions) Logger {
	hOpts := new(slog.HandlerOptions)
	for i := range opts {
		opts[i](hOpts)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, hOpts))

	return &Log{logger}
}

func (l *Log) WithField(key string, val string) Logger {
	return &Log{l.original.With(slog.Any(key, val))}
}

func (l *Log) WithFields(fields map[string]any) Logger {
	attrs := make([]any, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}

	return &Log{l.original.With(attrs...)}
}

func (l *Log) WithError(err error) Logger {
	return &Log{l.original.With(slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	})}
}

func (l *Log) Error(msg string) {
	l.original.Error(msg)
}

func (l *Log) Info(msg string) {
	l.original.Info(msg)
}

func (l *Log) Debug(msg string) {
	l.original.Debug(msg)
}

func (l *Log) Warn(msg string) {
	l.original.Warn(msg)
}
