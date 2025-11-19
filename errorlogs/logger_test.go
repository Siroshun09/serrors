package errorlogs

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
	"go.uber.org/mock/gomock"
)

func TestNewLogger(t *testing.T) {
	dedicated := logs.NewStdoutLogger(true)
	expect := logger{dedicated: dedicated}
	actual := NewLogger(dedicated)
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("expect: %+v, actual: %+v", expect, actual)
	}
}

func TestLogger_Debug(t *testing.T) {
	tests := []struct {
		name   string
		opt    LoggerOption
		msg    string
		expect func(ctx context.Context, mock *logmock.MockLogger)
	}{
		{
			name: "empty msg",
			opt:  LoggerOption{},
			msg:  "",
			expect: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Debug(ctx, "")
			},
		},
		{
			name: "not empty msg",
			opt:  LoggerOption{},
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

			l := logger{dedicated: mockLogger, opt: tt.opt}
			l.Debug(ctx, tt.msg)
		})
	}
}

func TestLogger_Info(t *testing.T) {
	tests := []struct {
		name   string
		opt    LoggerOption
		msg    string
		expect func(ctx context.Context, mock *logmock.MockLogger)
	}{
		{
			name: "empty msg",
			opt:  LoggerOption{},
			msg:  "",
			expect: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Info(ctx, "")
			},
		},
		{
			name: "not empty msg",
			opt:  LoggerOption{},
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

			l := logger{dedicated: mockLogger, opt: tt.opt}
			l.Info(ctx, tt.msg)
		})
	}
}

func TestLogger_Warn(t *testing.T) {
	tests := []struct {
		name   string
		opt    LoggerOption
		err    error
		expect func(ctx context.Context, err error, mock *logmock.MockLogger)
	}{
		{
			name: "empty error msg",
			opt:  LoggerOption{},
			err:  errors.New(""),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Warn(ctx, errors.New(""))
			},
		},
		{
			name: "not empty error msg",
			opt:  LoggerOption{},
			err:  errors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Warn(ctx, errors.New("test"))
			},
		},
		{
			name: "stacktrace attached / PrintStackTraceOnWarn = true",
			opt: LoggerOption{
				PrintStackTraceOnWarn: true,
			},
			err: serrors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Warn(ctx, err)
				mock.EXPECT().Debug(ctx, fmt.Sprintf(stackTraceLogFormat, serrors.GetStackTrace(err)))
			},
		},
		{
			name: "stacktrace attached / PrintStackTraceOnWarn = false",
			opt: LoggerOption{
				PrintStackTraceOnWarn: false,
			},
			err: serrors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Warn(ctx, err)
			},
		},
		{
			name: "stacktrace not attached / PrintStackTraceOnWarn = true / PrintCurrentStackTraceIfNotAttached = true",
			opt: LoggerOption{
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
			opt: LoggerOption{
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
			opt: LoggerOption{
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

			l := logger{dedicated: mockLogger, opt: tt.opt}
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
				mockLogger.EXPECT().Debug(ctx, fmt.Sprintf(stackTraceLogFormat, serrors.GetStackTrace(serr1)))
				mockLogger.EXPECT().Debug(ctx, fmt.Sprintf(stackTraceLogFormat, serrors.GetStackTrace(serr2)))
				mockLogger.EXPECT().Debug(ctx, fmt.Sprintf(stackTraceLogFormat, serrors.GetStackTrace(serr3)))
			}

			l := logger{dedicated: mockLogger, opt: LoggerOption{PrintStackTraceOnWarn: printStackTraceOnWarn}}
			l.Warn(ctx, err)
		})
	}

}

func TestLogger_Warnf(t *testing.T) {
	tests := []struct {
		name   string
		opt    LoggerOption
		format string
		arg    any
		expect func(ctx context.Context, mock *logmock.MockLogger)
	}{
		{
			name:   "empty msg",
			opt:    LoggerOption{},
			format: "%s",
			arg:    "",
			expect: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Warnf(ctx, "%s", "")
			},
		},
		{
			name:   "not empty msg",
			opt:    LoggerOption{},
			format: "%s",
			arg:    "test",
			expect: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Warnf(ctx, "%s", "test")
			},
		},
		{
			name:   "print current stacktrace / PrintStackTraceOnWarn = true / PrintCurrentStackTraceIfNotAttached = true",
			opt:    LoggerOption{PrintStackTraceOnWarn: true, PrintCurrentStackTraceIfNotAttached: true},
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

			l := logger{dedicated: mockLogger, opt: tt.opt}
			l.Warnf(ctx, tt.format, tt.arg)
		})
	}
}

func TestLogger_Error(t *testing.T) {
	tests := []struct {
		name   string
		opt    LoggerOption
		err    error
		expect func(ctx context.Context, err error, mock *logmock.MockLogger)
	}{
		{
			name: "empty error msg",
			opt:  LoggerOption{},
			err:  errors.New(""),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Error(ctx, errors.New(""))
			},
		},
		{
			name: "not empty error msg",
			opt:  LoggerOption{},
			err:  errors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Error(ctx, errors.New("test"))
			},
		},
		{
			name: "stacktrace attached",
			opt:  LoggerOption{},
			err:  serrors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Error(ctx, err)
				mock.EXPECT().Debug(ctx, fmt.Sprintf(stackTraceLogFormat, serrors.GetStackTrace(err)))
			},
		},
		{
			name: "stacktrace not attached / PrintCurrentStackTraceIfNotAttached = true",
			opt:  LoggerOption{PrintCurrentStackTraceIfNotAttached: true},
			err:  errors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Error(ctx, err)
				var stringType = reflect.TypeOf((*string)(nil)).Elem()
				mock.EXPECT().Debug(ctx, gomock.AssignableToTypeOf(stringType)) // print current stacktrace
			},
		},
		{
			name: "stacktrace not attached / PrintCurrentStackTraceIfNotAttached = false",
			opt:  LoggerOption{PrintCurrentStackTraceIfNotAttached: false},
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

			l := logger{dedicated: mockLogger, opt: tt.opt}
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
		mockLogger.EXPECT().Debug(ctx, fmt.Sprintf(stackTraceLogFormat, serrors.GetStackTrace(serr1)))
		mockLogger.EXPECT().Debug(ctx, fmt.Sprintf(stackTraceLogFormat, serrors.GetStackTrace(serr2)))
		mockLogger.EXPECT().Debug(ctx, fmt.Sprintf(stackTraceLogFormat, serrors.GetStackTrace(serr3)))

		l := logger{dedicated: mockLogger, opt: LoggerOption{}}
		l.Error(ctx, err)
	})
}

func TestLogger_Errorf(t *testing.T) {
	tests := []struct {
		name   string
		opt    LoggerOption
		format string
		arg    any
		expect func(ctx context.Context, mock *logmock.MockLogger)
	}{
		{
			name:   "empty msg",
			opt:    LoggerOption{},
			format: "%s",
			arg:    "",
			expect: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Errorf(ctx, "%s", "")
			},
		},
		{
			name:   "not empty msg",
			opt:    LoggerOption{},
			format: "%s",
			arg:    "test",
			expect: func(ctx context.Context, mock *logmock.MockLogger) {
				mock.EXPECT().Errorf(ctx, "%s", "test")
			},
		},
		{
			name:   "print current stacktrace / PrintCurrentStackTraceIfNotAttached = true",
			opt:    LoggerOption{PrintCurrentStackTraceIfNotAttached: true},
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

			l := logger{dedicated: mockLogger, opt: tt.opt}
			l.Errorf(ctx, tt.format, tt.arg)
		})
	}
}

func TestLogger_printStackTrace(t *testing.T) {
	tests := []struct {
		name   string
		opt    LoggerOption
		err    error
		expect func(ctx context.Context, err error, mock *logmock.MockLogger)
	}{
		{
			name: "nil",
			opt:  LoggerOption{},
			err:  nil,
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				// expect nothing to be called
			},
		},
		{
			name: "stacktrace attached error",
			opt:  LoggerOption{},
			err:  serrors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Debug(ctx, fmt.Sprintf(stackTraceLogFormat, serrors.GetStackTrace(err)))
			},
		},
		{
			name: "stacktrace not attached error",
			opt:  LoggerOption{},
			err:  errors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				// expect nothing to be called
			},
		},
		{
			name: "stacktrace not attached error / PrintCurrentStackTraceIfNotAttached = true",
			opt: LoggerOption{
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
			opt: LoggerOption{
				StackTraceLogLevel: StackTraceLogLevelInfo,
			},
			err: serrors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Info(ctx, fmt.Sprintf(stackTraceLogFormat, serrors.GetStackTrace(err)))
			},
		},
		{
			name: "stacktrace attached error / log level: warn",
			opt: LoggerOption{
				StackTraceLogLevel: StackTraceLogLevelWarn,
			},
			err: serrors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Warnf(ctx, stackTraceLogFormat, serrors.GetStackTrace(err))
			},
		},
		{
			name: "stacktrace attached error / log level: error",
			opt: LoggerOption{
				StackTraceLogLevel: StackTraceLogLevelError,
			},
			err: serrors.New("test"),
			expect: func(ctx context.Context, err error, mock *logmock.MockLogger) {
				mock.EXPECT().Errorf(ctx, stackTraceLogFormat, serrors.GetStackTrace(err))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			mockLogger := logmock.NewMockLogger(gomock.NewController(t))

			tt.expect(ctx, tt.err, mockLogger)

			l := logger{dedicated: mockLogger, opt: tt.opt}
			l.printStackTraces(ctx, tt.err)
		})
	}
}
