package logger

import (
	"context"
	"fmt"
	"go-scheduler/config"
	"os"
	"runtime"

	"github.com/rs/zerolog"
)

// logger is the global logger for this repository.
var logger zerolog.Logger

// InitZerolog initializes the zerolog logger.
func InitZerolog(config *config.Config) zerolog.Logger {
	logger = zerolog.New(os.Stderr).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()
	logger.Info().Msg("Logger initialized")
	if config.AppEnv == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return logger
}

// Info logs a message at info level.
func Info(message string) {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		logger.Info().Caller(1).Msg(fmt.Sprintf("[%s]%s", details.Name(), message))
	} else {
		logger.Info().Msg(message)
	}
}

// Warn logs a message at warn level.
func Warn(message string) {
	logger.Warn().Caller(1).Msg(message)
}

// Error logs a message at error level.
func Error(message string) {
	logger.Error().Caller(1).Msg(message)
}

// Debug logs a message at debug level.
func Debug(message string) {
	logger.Debug().Caller(1).Msg(message)
}

// Infof logs a message at info level with format.
func Infof(format string, args ...interface{}) {
	logger.Info().Msgf(format, args...)
}

// Warnf logs a message at warn level with format.
func Warnf(format string, args ...interface{}) {
	logger.Warn().Caller(1).Msgf(format, args...)
}

// Errorf logs a message at error level with format.
func Errorf(format string, args ...interface{}) {
	logger.Error().Caller(1).Msgf(format, args...)
}

// Debugf logs a message at debug level with format.
func Debugf(format string, args ...interface{}) {
	logger.Debug().Caller(1).Msgf(format, args...)
}

// InfoCtx logs a message at info level with context.
func InfoCtx(ctx context.Context, message string) {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		logger.Info().Interface("request_id", ctx.Value(config.RequestIDKey)).Caller(1).Msg(fmt.Sprintf("[%s]%s", details.Name(), message))
	} else {
		logger.Info().Interface("request_id", ctx.Value(config.RequestIDKey)).Msg(fmt.Sprintf("[%s]%s", ctx.Value(config.RequestIDKey), message))
	}
}

// WarnCtx logs a message at warn level with context.
func WarnCtx(ctx context.Context, message string) {
	logger.Warn().Interface("request_id", ctx.Value(config.RequestIDKey)).Caller(1).Msg(fmt.Sprintf("[%s]%s", ctx.Value(config.RequestIDKey), message))
}

// ErrorCtx logs a message at error level with context.
func ErrorCtx(ctx context.Context, message string) {
	logger.Error().Interface("request_id", ctx.Value(config.RequestIDKey)).Caller(1).Msg(fmt.Sprintf("[%s]%s", ctx.Value(config.RequestIDKey), message))
}

// DebugCtx logs a message at debug level with context.
func DebugCtx(ctx context.Context, message string) {
	logger.Debug().Interface("request_id", ctx.Value(config.RequestIDKey)).Caller(1).Msg(fmt.Sprintf("[%s]%s", ctx.Value(config.RequestIDKey), message))
}

// InfofCtx logs a message at info level with format with context.
func InfofCtx(ctx context.Context, format string, args ...interface{}) {
	logger.Info().Interface("request_id", ctx.Value(config.RequestIDKey)).Msgf(format, args...)
}

// WarnfCtx logs a message at warn level with format with context.
func WarnfCtx(ctx context.Context, format string, args ...interface{}) {
	logger.Warn().Interface("request_id", ctx.Value(config.RequestIDKey)).Caller(1).Msgf(format, args...)
}

// ErrorfCtx logs a message at error level with format with context.
func ErrorfCtx(ctx context.Context, format string, args ...interface{}) {
	logger.Error().Interface("request_id", ctx.Value(config.RequestIDKey)).Caller(1).Msgf(format, args...)
}

// DebugfCtx logs a message at debug level with format.
func DebugfCtx(ctx context.Context, format string, args ...interface{}) {
	logger.Debug().Interface("request_id", ctx.Value(config.RequestIDKey)).Caller(1).Msgf(format, args...)
}
