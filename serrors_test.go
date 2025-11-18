package serrors

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func TestWithStackTrace(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantNil     bool
		wantSame    bool
		wantWrapped bool
	}{
		{
			name:    "nil",
			err:     nil,
			wantNil: true,
		},
		{
			name:        "want wrapped err",
			err:         errors.New("test"),
			wantWrapped: true,
		},
		{
			name:     "want same",
			err:      New("test"),
			wantSame: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithStackTrace(tt.err)
			switch {
			case tt.wantNil:
				if got != nil {
					t.Errorf("want nil, got %v", got)
				}
			case tt.wantSame:
				{
					if !reflect.DeepEqual(got, tt.err) {
						t.Errorf("want %v, got %v", tt.err, got)
					}
				}
			case tt.wantWrapped:
				{
					if !errors.Is(got, tt.err) {
						t.Errorf("want %v, got %v", tt.err, got)
					}
				}
			default:
				t.Errorf("invalid test case")
			}
		})
	}
}

func Test_getStackTraceError(t *testing.T) {
	wrapped := New("wrapped")
	tests := []struct {
		name        string
		err         error
		wantSame    bool
		wantWrapped bool
		wantNil     bool
	}{
		{
			name:     "get same (New)",
			err:      New("test"),
			wantSame: true,
		},
		{
			name:     "get same (Errorf)",
			err:      Errorf("%s", "test"),
			wantSame: true,
		},
		{
			name:     "get same (WithStackTrace)",
			err:      WithStackTrace(errors.New("test")),
			wantSame: true,
		},
		{
			name:        "get wrapped (WithStackTrace with wrapped)",
			err:         WithStackTrace(wrapped),
			wantWrapped: true,
		},
		{
			name:        "get wrapped (WithStackTrace with fmt.Errorf with wrapped)",
			err:         WithStackTrace(fmt.Errorf("%w", wrapped)),
			wantWrapped: true,
		},
		{
			name:        "get wrapped (New)",
			err:         fmt.Errorf("%w", wrapped),
			wantWrapped: true,
		},
		{
			name:    "get nil (errors.New)",
			err:     errors.New("test"),
			wantNil: true,
		},
		{
			name:    "get nil (fmt.Errorf)",
			err:     fmt.Errorf("%s", "test"),
			wantNil: true,
		},
		{
			name:    "get nil (nil)",
			err:     nil,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStackTraceError(tt.err)

			switch {
			case tt.wantSame:
				if !reflect.DeepEqual(got, tt.err) {
					t.Errorf("getStackTraceError() = %v, want %v", got, tt.err)
				}
			case tt.wantWrapped:
				if !reflect.DeepEqual(got, wrapped) {
					t.Errorf("getStackTraceError() = %v, want %v", got, wrapped)
				}
			case tt.wantNil:
				if got != nil {
					t.Errorf("getStackTraceError() = %v, want %v", got, nil)
				}
			default:
				t.Errorf("invalid test case")
			}
		})
	}
}

func TestGetStackTrace(t *testing.T) {
	stackTrace := StackTrace{}
	tests := []struct {
		name        string
		err         error
		want        StackTrace
		wantCreated bool
	}{
		{
			name: "expect nil",
			err:  nil,
			want: nil,
		},
		{
			name: "expect stackTrace",
			err:  &stackTraceError{err: errors.New("test"), stackTrace: stackTrace},
			want: stackTrace,
		},
		{
			name: "expect stackTrace (wrapped)",
			err:  fmt.Errorf("%w", &stackTraceError{err: errors.New("test"), stackTrace: stackTrace}),
			want: stackTrace,
		},
		{
			name:        "expect stack trace created by GetStackTrace",
			err:         errors.New("test"),
			wantCreated: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantCreated {
				if got := GetStackTrace(tt.err); got == nil {
					t.Errorf("GetStackTrace() = %v, want created stack trace", got)
				}
			} else if got := GetStackTrace(tt.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStackTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAttachedStackTrace(t *testing.T) {
	stackTrace := StackTrace{}
	tests := []struct {
		name        string
		err         error
		want        StackTrace
		wantCreated bool
	}{
		{
			name: "expect nil and false (err is nil)",
			err:  nil,
			want: nil,
		},
		{
			name: "expect stackTrace and true",
			err:  &stackTraceError{err: errors.New("test"), stackTrace: stackTrace},
			want: stackTrace,
		},
		{
			name: "expect stackTrace (wrapped) and true",
			err:  fmt.Errorf("%w", &stackTraceError{err: errors.New("test"), stackTrace: stackTrace}),
			want: stackTrace,
		},
		{
			name:        "expect nil and false (err is not stackTraceError)",
			err:         errors.New("test"),
			wantCreated: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := GetAttachedStackTrace(tt.err)

			var fail bool
			if tt.want == nil {
				fail = ok || got != nil
			} else {
				fail = !ok || !reflect.DeepEqual(got, tt.want)
			}

			if fail {
				t.Errorf("GetAttachedStackTrace() = (%v, %v), want (%v, %v)", got, ok, tt.want, tt.want != nil)
			}
		})
	}
}

func TestGetCurrentStackTrace(t *testing.T) {
	stackTrace := GetCurrentStackTrace()
	if len(stackTrace) == 0 {
		t.Errorf("GetCurrentStackTrace() = %v, want not empty", stackTrace)
		return
	}
}

func TestFuncInfo_String(t *testing.T) {

	tests := []struct {
		name     string
		funcInfo FuncInfo
		want     string
	}{
		{
			name: "success",
			funcInfo: FuncInfo{
				Name: "Test",
				File: "test.go",
				Line: 10,
			},
			want: "Test (test.go:10)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.funcInfo.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStackTrace_String(t *testing.T) {
	tests := []struct {
		name       string
		stackTrace StackTrace
		want       string
	}{
		{
			name:       "empty",
			stackTrace: StackTrace{},
			want:       "",
		},
		{
			name: "1 FuncInfo",
			stackTrace: StackTrace{
				{
					Name: "Test",
					File: "test.go",
					Line: 1,
				},
			},
			want: "Test (test.go:1)",
		},
		{
			name: "2 FuncInfo",
			stackTrace: StackTrace{
				{
					Name: "Test1",
					File: "test.go",
					Line: 1,
				},
				{
					Name: "Test2",
					File: "test.go",
					Line: 2,
				},
			},
			want: "Test1 (test.go:1)\nTest2 (test.go:2)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.stackTrace.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newFuncInfo(t *testing.T) {
	pc, _, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime error, can't get caller")
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		t.Fatalf("runtime error, can't get Func")
	}
	tests := []struct {
		name        string
		pc          uintptr
		f           *runtime.Func
		wantUnknown bool
	}{
		{
			name:        "get unknown",
			pc:          0,
			f:           nil,
			wantUnknown: true,
		},
		{
			name:        "get actual",
			pc:          pc,
			f:           f,
			wantUnknown: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newFuncInfo(tt.pc, tt.f)
			if tt.wantUnknown {
				if got != unknownFuncInfo {
					t.Errorf("newFuncInfo() = %v, want %v", got, unknownFuncInfo)
				}
			} else {
				if got == unknownFuncInfo {
					t.Errorf("newFuncInfo() = %v, want actual func info", unknownFuncInfo)
				}
			}
		})
	}
}
