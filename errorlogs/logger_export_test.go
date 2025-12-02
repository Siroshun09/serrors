package errorlogs

import (
	"context"
	"reflect"

	"github.com/Siroshun09/logs"
)

func IsLoggerDedicatedBy(expect logs.Logger, actual logs.Logger) bool {
	wrapped, ok := actual.(logger)
	if !ok {
		return false
	}
	return reflect.DeepEqual(expect, wrapped.dedicated)
}

func CallPrintStackTrace(ctx context.Context, err error, target logs.Logger) {
	wrapped, ok := target.(logger)
	if !ok {
		panic("invalid logger")
	}
	wrapped.printStackTraces(ctx, err)
}

// GetStackTraceLogFormat exposes the internal stackTraceLogFormat for external tests.
func GetStackTraceLogFormat() string {
	return stackTraceLogFormat
}
