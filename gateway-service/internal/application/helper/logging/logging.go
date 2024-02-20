package logging

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxLoggerKey struct{}

func ContextWithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, l)
}

func LoggerFromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxLoggerKey{}).(*zap.Logger); ok {
		return l
	}

	return zap.L()
}

func ZapFromEnv() (logger *zap.Logger, syncFunc func(), err error) {
	type syncable interface {
		Sync() error
	}

	var syncables []syncable

	logCfg := zap.NewProductionConfig()
	logCfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logCfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)

	logger, err = logCfg.Build()
	if err != nil {
		return nil, nil, err
	}

	syncables = append(syncables, logger)
	zap.ReplaceGlobals(logger)
	return logger, func() {
		for _, s := range syncables {
			s.Sync()
		}
	}, nil
}
