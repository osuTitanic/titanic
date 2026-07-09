package database

import (
	"context"
	"log/slog"
	"time"

	gormLogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	gormLogger.Interface
	logger *slog.Logger
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		logger: slog.Default().With("component", "database"),
	}
}

func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...any) {
	l.logger.Info(msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	l.logger.Warn(msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...any) {
	l.logger.Error(msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err == nil {
		l.logger.Debug("Trace", "elapsed", elapsed, "rows", rows, "sql", sql)
		return
	}

	if err != gormLogger.ErrRecordNotFound {
		l.logger.Error("Trace", "error", err, "elapsed", elapsed, "rows", rows, "sql", sql)
	} else {
		l.logger.Debug("Trace", "error", err, "elapsed", elapsed, "rows", rows, "sql", sql)
	}
}
