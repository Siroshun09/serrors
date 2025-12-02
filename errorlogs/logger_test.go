package errorlogs_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/Siroshun09/logs"
	"github.com/Siroshun09/logs/logmock"
	"github.com/Siroshun09/serrors"
	"github.com/Siroshun09/serrors/errorlogs"
	"go.uber.org/mock/gomock"
)

func TestNewLogger(t *testing.T) {
	dedicated := logs.NewStdoutLogger(true)
	actual := errorlogs.NewLogger(dedicated)
	if !errorlogs.IsLoggerDedicatedBy(dedicated, actual) {
		t.Errorf("expect: %+v, actual: %+v", dedicated, actual)
	}
}

func TestLogger_Debug(t *testing.T) {
	tests := []struct {
		name   string
		opt    errorlogs.LoggerOption
		msg    string
		expect func(ctx context.Context, mock *logmock.MockLogger)
	}{
		{
			name: "empty msg",
			opt:  errorlogs.LoggerOption{},
			msg:  "",
			expect: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Debug(ctx, "")
			},
		},
		{
			name: "not empty msg",
			opt:  errorlogs.LoggerOption{},
			msg:  "test",
			expect: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Debug(ctx, "test")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			mockLogger := logmock.NewMockLogger(gomock.NewController(t))

			tt.expect(ctx, mockLogger)

			l := errorlogs.NewLoggerWithOption(mockLogger, tt.opt)
			l.Debug(ctx, tt.msg)
		})
	}
}

func TestLogger_Info(t *testing.T) {
	tests := []struct {
		name   string
		opt    errorlogs.LoggerOption
		msg    string
		expect func(ctx context.Context, mock *logmock.MockLogger)
	}{
		{
			name: "empty msg",
			opt:  errorlogs.LoggerOption{},
			msg:  "",
			expect: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Info(ctx, "")
			},
		},
		{
			name: "not empty msg",
			opt:  errorlogs.LoggerOption{},
			msg:  "test",
			expect: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Info(ctx, "test")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			mockLogger := logmock.NewMockLogger(gomock.NewController(t))

			tt.expect(ctx, mockLogger)

			l := errorlogs.NewLoggerWithOption(mockLogger, tt.opt)
			l.Info(ctx, tt.msg)
		})
	}
}

func TestLogger_Warn(t *testing.T) {
	tests := []struct {
		name   string
		opt    errorlogs.LoggerOption
		err    error
		expect func(ctx context.Context, err error, mock *logmock.MockLogger)
	}{
		{
			name: "empty error msg",
			opt:  errorlogs.LoggerOption{},
			err:  errors.New(""),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Warn(ctx, errors.New(""))
			},
		},
		{
			name: "not empty error msg",
			opt:  errorlogs.LoggerOption{},
			err:  errors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Warn(ctx, errors.New("test"))
			},
		},
		{
			name: "stacktrace attached / PrintStackTraceOnWarn = true",
			opt: errorlogs.LoggerOption{
				PrintStackTraceOnWarn: true,
			},
			err: serrors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Warn(ctx, err)
				mock.EXPECT().Debug(ctx, fmt.Sprintf(errorlogs.GetStackTraceLogFormat(), serrors.GetStackTrace(err)))
			},
		},
		{
			name: "stacktrace attached / PrintStackTraceOnWarn = false",
			opt: errorlogs.LoggerOption{
				PrintStackTraceOnWarn: false,
			},
			err: serrors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Warn(ctx, err)
			},
		},
		{
			name: "stacktrace not attached / PrintStackTraceOnWarn = true / PrintCurrentStackTraceIfNotAttached = true",
			opt: errorlogs.LoggerOption{
				PrintStackTraceOnWarn:               true,
				PrintCurrentStackTraceIfNotAttached: true,
			},
			err: errors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Warn(ctx, err)
				var stringType = reflect.TypeOf((*string)(nil)).Elem()
				mock.EXPECT().Debug(ctx, gomock.AssignableToTypeOf(stringType)) // print current stacktrace
			},
		},
		{
			name: "stacktrace not attached / PrintStackTraceOnWarn = false / PrintCurrentStackTraceIfNotAttached = true",
			opt: errorlogs.LoggerOption{
				PrintStackTraceOnWarn:               false,
				PrintCurrentStackTraceIfNotAttached: true,
			},
			err: errors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Warn(ctx, err)
			},
		},
		{
			name: "stacktrace not attached / PrintStackTraceOnWarn = true / PrintCurrentStackTraceIfNotAttached = false",
			opt: errorlogs.LoggerOption{
				PrintStackTraceOnWarn:               true,
				PrintCurrentStackTraceIfNotAttached: false,
			},
			err: errors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Warn(ctx, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			mockLogger := logmock.NewMockLogger(gomock.NewController(t))

			tt.expect(ctx, tt.err, mockLogger)

			l := errorlogs.NewLoggerWithOption(mockLogger, tt.opt)
			l.Warn(ctx, tt.err)
		})
	}

	for _, printStackTraceOnWarn := range []bool{true, false} {
		t.Run("multiple stacktraces / PrintStackTraceOnWarn = "+strconv.FormatBool(printStackTraceOnWarn), func(t *testing.T) {
			serr1 := serrors.New("test1")
			serr2 := serrors.New("test2")
			serr3 := serrors.New("test3")
			err := errors.Join(serr1, serr2, serr3)

			ctx := t.Context()
			mockLogger := logmock.NewMockLogger(gomock.NewController(t))

			mockLogger.EXPECT().Warn(ctx, err)
			if printStackTraceOnWarn {
				mockLogger.EXPECT().Debug(ctx, fmt.Sprintf(errorlogs.GetStackTraceLogFormat(), serrors.GetStackTrace(serr1)))
				mockLogger.EXPECT().Debug(ctx, fmt.Sprintf(errorlogs.GetStackTraceLogFormat(), serrors.GetStackTrace(serr2)))
				mockLogger.EXPECT().Debug(ctx, fmt.Sprintf(errorlogs.GetStackTraceLogFormat(), serrors.GetStackTrace(serr3)))
			}

			l := errorlogs.NewLoggerWithOption(mockLogger, errorlogs.LoggerOption{PrintStackTraceOnWarn: printStackTraceOnWarn})
			l.Warn(ctx, err)
		})
	}

}

func TestLogger_Warnf(t *testing.T) {
	tests := []struct {
		name   string
		opt    errorlogs.LoggerOption
		format string
		arg    any
		expect func(ctx context.Context, mock *logmock.MockLogger)
	}{
		{
			name:   "empty msg",
			opt:    errorlogs.LoggerOption{},
			format: "%s",
			arg:    "",
			expect: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Warnf(ctx, "%s", "")
			},
		},
		{
			name:   "not empty msg",
			opt:    errorlogs.LoggerOption{},
			format: "%s",
			arg:    "test",
			expect: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Warnf(ctx, "%s", "test")
			},
		},
		{
			name:   "print current stacktrace / PrintStackTraceOnWarn = true / PrintCurrentStackTraceIfNotAttached = true",
			opt:    errorlogs.LoggerOption{PrintStackTraceOnWarn: true, PrintCurrentStackTraceIfNotAttached: true},
			format: "%s",
			arg:    "test",
			expect: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Warnf(ctx, "%s", "test")
				var stringType = reflect.TypeOf((*string)(nil)).Elem()
				mock.EXPECT().Debug(ctx, gomock.AssignableToTypeOf(stringType)) // print current stacktrace
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			mockLogger := logmock.NewMockLogger(gomock.NewController(t))

			tt.expect(ctx, mockLogger)

			l := errorlogs.NewLoggerWithOption(mockLogger, tt.opt)
			l.Warnf(ctx, tt.format, tt.arg)
		})
	}
}

func TestLogger_Error(t *testing.T) {
	tests := []struct {
		name   string
		opt    errorlogs.LoggerOption
		err    error
		expect func(ctx context.Context, err error, mock *logmock.MockLogger)
	}{
		{
			name: "empty error msg",
			opt:  errorlogs.LoggerOption{},
			err:  errors.New(""),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Error(ctx, errors.New(""))
			},
		},
		{
			name: "not empty error msg",
			opt:  errorlogs.LoggerOption{},
			err:  errors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Error(ctx, errors.New("test"))
			},
		},
		{
			name: "stacktrace attached",
			opt:  errorlogs.LoggerOption{},
			err:  serrors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Error(ctx, err)
				mock.EXPECT().Debug(ctx, fmt.Sprintf(errorlogs.GetStackTraceLogFormat(), serrors.GetStackTrace(err)))
			},
		},
		{
			name: "stacktrace not attached / PrintCurrentStackTraceIfNotAttached = true",
			opt:  errorlogs.LoggerOption{PrintCurrentStackTraceIfNotAttached: true},
			err:  errors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Error(ctx, err)
				var stringType = reflect.TypeOf((*string)(nil)).Elem()
				mock.EXPECT().Debug(ctx, gomock.AssignableToTypeOf(stringType)) // print current stacktrace
			},
		},
		{
			name: "stacktrace not attached / PrintCurrentStackTraceIfNotAttached = false",
			opt:  errorlogs.LoggerOption{PrintCurrentStackTraceIfNotAttached: false},
			err:  errors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Error(ctx, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			mockLogger := logmock.NewMockLogger(gomock.NewController(t))

			tt.expect(ctx, tt.err, mockLogger)

			l := errorlogs.NewLoggerWithOption(mockLogger, tt.opt)
			l.Error(ctx, tt.err)
		})
	}

	t.Run("multiple stacktraces", func(t *testing.T) {
		serr1 := serrors.New("test1")
		serr2 := serrors.New("test2")
		serr3 := serrors.New("test3")
		err := errors.Join(serr1, serr2, serr3)

		ctx := t.Context()
		mockLogger := logmock.NewMockLogger(gomock.NewController(t))

		mockLogger.EXPECT().Error(ctx, err)
		mockLogger.EXPECT().Debug(ctx, fmt.Sprintf(errorlogs.GetStackTraceLogFormat(), serrors.GetStackTrace(serr1)))
		mockLogger.EXPECT().Debug(ctx, fmt.Sprintf(errorlogs.GetStackTraceLogFormat(), serrors.GetStackTrace(serr2)))
		mockLogger.EXPECT().Debug(ctx, fmt.Sprintf(errorlogs.GetStackTraceLogFormat(), serrors.GetStackTrace(serr3)))

		l := errorlogs.NewLoggerWithOption(mockLogger, errorlogs.LoggerOption{})
		l.Error(ctx, err)
	})
}

func TestLogger_Errorf(t *testing.T) {
	tests := []struct {
		name   string
		opt    errorlogs.LoggerOption
		format string
		arg    any
		expect func(ctx context.Context, mock *logmock.MockLogger)
	}{
		{
			name:   "empty msg",
			opt:    errorlogs.LoggerOption{},
			format: "%s",
			arg:    "",
			expect: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Errorf(ctx, "%s", "")
			},
		},
		{
			name:   "not empty msg",
			opt:    errorlogs.LoggerOption{},
			format: "%s",
			arg:    "test",
			expect: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Errorf(ctx, "%s", "test")
			},
		},
		{
			name:   "print current stacktrace / PrintCurrentStackTraceIfNotAttached = true",
			opt:    errorlogs.LoggerOption{PrintCurrentStackTraceIfNotAttached: true},
			format: "%s",
			arg:    "test",
			expect: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Errorf(ctx, "%s", "test")
				var stringType = reflect.TypeOf((*string)(nil)).Elem()
				mock.EXPECT().Debug(ctx, gomock.AssignableToTypeOf(stringType)) // print current stacktrace
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			mockLogger := logmock.NewMockLogger(gomock.NewController(t))

			tt.expect(ctx, mockLogger)

			l := errorlogs.NewLoggerWithOption(mockLogger, tt.opt)
			l.Errorf(ctx, tt.format, tt.arg)
		})
	}
}

func TestLogger_printStackTrace(t *testing.T) {
	tests := []struct {
		name   string
		opt    errorlogs.LoggerOption
		err    error
		expect func(ctx context.Context, err error, mock *logmock.MockLogger)
	}{
		{
			name: "nil",
			opt:  errorlogs.LoggerOption{},
			err:  nil,
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				// expect nothing to be called
			},
		},
		{
			name: "stacktrace attached error",
			opt:  errorlogs.LoggerOption{},
			err:  serrors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Debug(ctx, fmt.Sprintf(errorlogs.GetStackTraceLogFormat(), serrors.GetStackTrace(err)))
			},
		},
		{
			name: "stacktrace not attached error",
			opt:  errorlogs.LoggerOption{},
			err:  errors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				// expect nothing to be called
			},
		},
		{
			name: "stacktrace not attached error / PrintCurrentStackTraceIfNotAttached = true",
			opt: errorlogs.LoggerOption{
				PrintCurrentStackTraceIfNotAttached: true,
			},
			err: errors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				var stringType = reflect.TypeOf((*string)(nil)).Elem()
				mock.EXPECT().Debug(ctx, gomock.AssignableToTypeOf(stringType)) // print current stacktrace
			},
		},
		{
			name: "stacktrace attached error / log level: info",
			opt: errorlogs.LoggerOption{
				StackTraceLogLevel: errorlogs.StackTraceLogLevelInfo,
			},
			err: serrors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Info(ctx, fmt.Sprintf(errorlogs.GetStackTraceLogFormat(), serrors.GetStackTrace(err)))
			},
		},
		{
			name: "stacktrace attached error / log level: warn",
			opt: errorlogs.LoggerOption{
				StackTraceLogLevel: errorlogs.StackTraceLogLevelWarn,
			},
			err: serrors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Warnf(ctx, errorlogs.GetStackTraceLogFormat(), serrors.GetStackTrace(err))
			},
		},
		{
			name: "stacktrace attached error / log level: error",
			opt: errorlogs.LoggerOption{
				StackTraceLogLevel: errorlogs.StackTraceLogLevelError,
			},
			err: serrors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Errorf(ctx, errorlogs.GetStackTraceLogFormat(), serrors.GetStackTrace(err))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			mockLogger := logmock.NewMockLogger(gomock.NewController(t))

			tt.expect(ctx, tt.err, mockLogger)

			l := errorlogs.NewLoggerWithOption(mockLogger, tt.opt)
			errorlogs.CallPrintStackTrace(ctx, tt.err, l)
		})
	}
}
