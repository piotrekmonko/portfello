package logz

import (
	"context"
	"fmt"
	"github.com/piotrekmonko/portfello/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerKeyValuePairsCtx struct{}

var (
	AtomicLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	version     = "dev"
	loggerCtx   = loggerKeyValuePairsCtx{}
)

type Logger interface {
	Named(name string) Logger
	With(args ...interface{}) Logger
	Debugw(ctx context.Context, msg string, keysAndValues ...interface{})
	Infow(ctx context.Context, msg string, keysAndValues ...interface{})
	Warnw(ctx context.Context, msg string, keysAndValues ...interface{})
	Errorw(ctx context.Context, err error, msg string, keysAndValues ...interface{}) error
}

func NewLogger(c *config.Logging) Logger {
	var cfg zap.Config

	if c.Format == "dev" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	if err := ParseFlag(c.Level); err != nil {
		panic(err)
	}

	cfg.Level.SetLevel(AtomicLevel.Level())
	l, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	return &Log{SugaredLogger: l.Sugar().With("v", version)}
}

var _ Logger = (*Log)(nil)

type Log struct {
	*zap.SugaredLogger
}

func (l *Log) Named(name string) Logger {
	return &Log{SugaredLogger: l.SugaredLogger.Named(name)}
}

func (l *Log) With(args ...interface{}) Logger {
	return &Log{SugaredLogger: l.SugaredLogger.With(args...)}
}

func (l *Log) Debugw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.With(FromCtx(ctx)...).Debugw(msg, keysAndValues...)
}

func (l *Log) Infow(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.With(FromCtx(ctx)...).Infow(msg, keysAndValues...)
}

func (l *Log) Warnw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.With(FromCtx(ctx)...).Warnw(msg, keysAndValues...)
}

func (l *Log) Errorw(ctx context.Context, err error, msg string, keysAndValues ...interface{}) error {
	if err == nil {
		return nil
	}

	l.SugaredLogger.With(zap.Error(err)).With(FromCtx(ctx)...).Errorw(msg, keysAndValues...)
	return fmt.Errorf("%s: %w", msg, err)
}

func WithCtx(ctx context.Context, keysAndValues ...any) context.Context {
	pairs, havePairs := ctx.Value(loggerCtx).([]any)
	if !havePairs {
		pairs = keysAndValues
	} else {
		pairs = append(pairs, keysAndValues...)
	}

	return context.WithValue(ctx, loggerCtx, pairs)
}

func FromCtx(ctx context.Context) []any {
	if ctx == nil {
		return []any{}
	}

	pairs, havePairs := ctx.Value(loggerCtx).([]any)
	if havePairs {
		return pairs
	}

	return []any{}
}

// ParseFlag sets logging level to "level", see zapcore.Level for possible values.
func ParseFlag(level string) error {
	al, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return fmt.Errorf("cannot parse log level: %w", err)
	}

	AtomicLevel = al
	return nil
}

// SetVer sets input as the value for "v" field (version) in log entries.
func SetVer(v string) {
	version = v
}
