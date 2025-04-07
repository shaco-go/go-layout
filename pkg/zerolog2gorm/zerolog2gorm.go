package zerolog2gorm

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type Logger struct {
	ZeroLog                   *zerolog.Logger
	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

func New(zapLogger *zerolog.Logger) Logger {
	return Logger{
		ZeroLog:                   zapLogger,
		LogLevel:                  logger.Warn,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
	}
}

func (l Logger) SetAsDefault() {
	logger.Default = l
}

func (l Logger) LogMode(level logger.LogLevel) logger.Interface {
	return Logger{
		ZeroLog:                   l.ZeroLog,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

func (l Logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Info {
		return
	}
	l.ZeroLog.Info().Ctx(ctx).Msgf(str, args...)
}

func (l Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Warn {
		return
	}
	l.ZeroLog.Warn().Ctx(ctx).Msgf(str, args...)
}

func (l Logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Error {
		return
	}
	l.ZeroLog.Error().Ctx(ctx).Msgf(str, args...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		l.ZeroLog.Error().
			Ctx(ctx).
			Err(err).
			Str("sql", sql).
			Msgf("elapsed %+v rows %+v", elapsed, rows)

	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		l.ZeroLog.Warn().
			Ctx(ctx).
			Str("sql", sql).
			Msgf("elapsed %+v rows %+v", elapsed, rows)

	case l.LogLevel >= logger.Info:
		sql, rows := fc()
		l.ZeroLog.Trace().
			Ctx(ctx).
			Str("sql", sql).
			Msgf("elapsed %+v rows %+v", elapsed, rows)
	}
}
