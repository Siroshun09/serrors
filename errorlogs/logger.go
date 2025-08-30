package errorlogs

import (
	"context"
	"fmt"

	"github.com/Siroshun09/logs"
	"github.com/Siroshun09/serrors"
)

// NewLogger creates a new logs.Logger.
func NewLogger(out logs.Logger) logs.Logger {
	return NewLoggerWithOption(out, LoggerOption{})
}

// NewLoggerWithOption creates a new logs.Logger with the given option.
func NewLoggerWithOption(out logs.Logger, opt LoggerOption) logs.Logger {
	return logger{
		dedicated: out,
		opt:       opt,
	}
}

// LoggerOption is the option for logger implementation.
type LoggerOption struct {
	// StackTraceLogLevel is the log level for stack trace.
	StackTraceLogLevel StackTraceLogLevel
	// PrintStackTraceOnWarn is whether to print stack trace on Warn.
	PrintStackTraceOnWarn bool
	// PrintCurrentStackTraceIfNotAttached is whether to print the current stack trace if the error does not have a stack trace.
	PrintCurrentStackTraceIfNotAttached bool
}

// StackTraceLogLevel is the log level for stack trace.
type StackTraceLogLevel int8

const (
	// StackTraceLogLevelDebug indicates that the stack trace will be logged as a debug level.
	StackTraceLogLevelDebug StackTraceLogLevel = iota
	// StackTraceLogLevelInfo indicates that the stack trace will be logged as an info level.
	StackTraceLogLevelInfo
	// StackTraceLogLevelWarn indicates that the stack trace will be logged as a warn level.
	StackTraceLogLevelWarn
	// StackTraceLogLevelError indicates that the stack trace will be logged as an error level.
	StackTraceLogLevelError
)

type logger struct {
	dedicated logs.Logger
	opt       LoggerOption
}

func (l logger) Debug(ctx context.Context, msg string) {
	l.dedicated.Debug(ctx, msg)
}

func (l logger) Info(ctx context.Context, msg string) {
	l.dedicated.Info(ctx, msg)
}

func (l logger) Warn(ctx context.Context, err error) {
	l.dedicated.Warn(ctx, err)
	if l.opt.PrintStackTraceOnWarn {
		l.printStackTrace(ctx, err)
	}
}

func (l logger) Warnf(ctx context.Context, format string, args ...any) {
	l.dedicated.Warnf(ctx, format, args...)
	if l.opt.PrintStackTraceOnWarn {
		l.printStackTrace(ctx, nil)
	}
}

func (l logger) Error(ctx context.Context, err error) {
	l.dedicated.Error(ctx, err)
	l.printStackTrace(ctx, err)
}

func (l logger) Errorf(ctx context.Context, format string, args ...any) {
	l.dedicated.Errorf(ctx, format, args...)
	l.printStackTrace(ctx, nil)
}

const stackTraceLogFormat = "stacktrace\n%s"

func (l logger) printStackTrace(ctx context.Context, err error) {
	stackTrace, exists := serrors.GetAttachedStackTrace(err)
	if !exists {
		if l.opt.PrintCurrentStackTraceIfNotAttached {
			stackTrace = serrors.GetCurrentStackTrace()
		} else {
			return
		}
	}

	switch l.opt.StackTraceLogLevel {
	case StackTraceLogLevelDebug:
		l.dedicated.Debug(ctx, fmt.Sprintf(stackTraceLogFormat, stackTrace))
	case StackTraceLogLevelInfo:
		l.dedicated.Info(ctx, fmt.Sprintf(stackTraceLogFormat, stackTrace))
	case StackTraceLogLevelWarn:
		l.dedicated.Warnf(ctx, stackTraceLogFormat, stackTrace)
	case StackTraceLogLevelError:
		l.dedicated.Errorf(ctx, stackTraceLogFormat, stackTrace)
	}
}
