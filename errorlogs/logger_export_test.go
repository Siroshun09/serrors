package errorlogs

import (
	"context"
	"reflect"

	"github.com/Siroshun09/logs"
	"github.com/Siroshun09/serrors"
)

func IsLoggerDedicatedBy(expect logs.Logger, actual logs.Logger) bool {
	return reflect.DeepEqual(expect, castLogger(actual).dedicated)
}

func GetLoggerOption(actual logs.Logger) LoggerOption {
	return castLogger(actual).opt
}

func CallPrintStackTraces(ctx context.Context, err error, target logs.Logger) {
	castLogger(target).printStackTraces(ctx, err)
}

func CallPrintStackTrace(ctx context.Context, target logs.Logger) {
	castLogger(target).printStackTrace(ctx, serrors.GetCurrentStackTrace())
}

// GetStackTraceLogFormat exposes the internal stackTraceLogFormat for external tests.
func GetStackTraceLogFormat() string {
	return stackTraceLogFormat
}

func NewNilLogger() logs.Logger {
	return (*logger)(nil)
}

func castLogger(l logs.Logger) *logger {
	casted, ok := l.(*logger)
	if !ok {
		panic("invalid logger")
	}
	return casted
}
