package serrors_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/Siroshun09/serrors/v2"
)

func TestFrame_String(t *testing.T) {
	tests := []struct {
		name  string
		frame serrors.Frame
		want  string
	}{
		{
			name: "standard",
			frame: serrors.Frame{
				Function: "main.main",
				File:     "main.go",
				Line:     10,
			},
			want: "main.main (main.go:10)",
		},
		{
			name: "empty function",
			frame: serrors.Frame{
				Function: "",
				File:     "main.go",
				Line:     10,
			},
			want: " (main.go:10)",
		},
		{
			name: "empty file",
			frame: serrors.Frame{
				Function: "main.main",
				File:     "",
				Line:     10,
			},
			want: "main.main (:10)",
		},
		{
			name: "zero line",
			frame: serrors.Frame{
				Function: "main.main",
				File:     "main.go",
				Line:     0,
			},
			want: "main.main (main.go:0)",
		},
		{
			name: "negative line",
			frame: serrors.Frame{
				Function: "main.main",
				File:     "main.go",
				Line:     -1,
			},
			want: "main.main (main.go:-1)",
		},
		{
			name: "all empty/zero",
			frame: serrors.Frame{
				Function: "",
				File:     "",
				Line:     0,
			},
			want: " (:0)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.frame.String(); got != tt.want {
				t.Errorf("Frame.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFrame_AppendText(t *testing.T) {
	tests := []struct {
		name   string
		frame  serrors.Frame
		prefix []byte
		want   []byte
	}{
		{
			name: "standard",
			frame: serrors.Frame{
				Function: "main.main",
				File:     "main.go",
				Line:     10,
			},
			prefix: nil,
			want:   []byte("main.main (main.go:10)"),
		},
		{
			name: "with prefix",
			frame: serrors.Frame{
				Function: "main.main",
				File:     "main.go",
				Line:     10,
			},
			prefix: []byte("prefix: "),
			want:   []byte("prefix: main.main (main.go:10)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.frame.AppendText(tt.prefix)
			if err != nil {
				t.Fatalf("Frame.AppendText() unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Frame.AppendText() = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}

func TestStackTrace_String(t *testing.T) {
	tests := []struct {
		name string
		st   serrors.StackTrace
		want string
	}{
		{
			name: "empty",
			st:   serrors.StackTrace{},
			want: "",
		},
		{
			name: "single",
			st: serrors.StackTrace{
				{Function: "f1", File: "file1.go", Line: 1},
			},
			want: "f1 (file1.go:1)",
		},
		{
			name: "multiple",
			st: serrors.StackTrace{
				{Function: "f1", File: "file1.go", Line: 1},
				{Function: "f2", File: "file2.go", Line: 2},
			},
			want: "f1 (file1.go:1)\nf2 (file2.go:2)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.st.String(); got != tt.want {
				t.Errorf("StackTrace.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStackTrace_AppendText(t *testing.T) {
	tests := []struct {
		name   string
		st     serrors.StackTrace
		prefix []byte
		want   []byte
	}{
		{
			name: "empty",
			st:   serrors.StackTrace{},
			want: nil,
		},
		{
			name: "single",
			st: serrors.StackTrace{
				{Function: "f1", File: "file1.go", Line: 1},
			},
			prefix: []byte("prefix: "),
			want:   []byte("prefix: f1 (file1.go:1)"),
		},
		{
			name: "multiple",
			st: serrors.StackTrace{
				{Function: "f1", File: "file1.go", Line: 1},
				{Function: "f2", File: "file2.go", Line: 2},
			},
			prefix: []byte("stack:\n"),
			want:   []byte("stack:\nf1 (file1.go:1)\nf2 (file2.go:2)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.st.AppendText(tt.prefix)
			if err != nil {
				t.Fatalf("StackTrace.AppendText() unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StackTrace.AppendText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCurrentStackTrace(t *testing.T) {
	stackTrace := serrors.GetCurrentStackTrace()
	if len(stackTrace) == 0 {
		t.Errorf("GetCurrentStackTrace() = %v, want not empty", stackTrace)
		return
	}
}

func Test_newStackTrace(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		st := serrors.NewStackTrace(1, 10)
		if len(st) == 0 {
			t.Fatal("expected non-empty stack trace")
		}

		if st[0].Function != "github.com/Siroshun09/serrors/v2_test.Test_newStackTrace.func1" {
			t.Errorf("could not find current function in stack trace: %v", st)
		}
	})

	t.Run("limit", func(t *testing.T) {
		st := serrors.NewStackTrace(0, 1)
		if len(st) != 1 {
			t.Errorf("expected stack trace length 1, got %d", len(st))
		}
	})
}

func Test_stackTraceError_Error(t *testing.T) {
	stackTrace := serrors.StackTrace{
		{Function: "f1", File: "file1.go", Line: 1},
		{Function: "f2", File: "file2.go", Line: 2},
	}

	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "wrap nil error",
			err:  serrors.NewStackTraceError(nil, stackTrace),
			want: "",
		},
		{
			name: "wrap non-nil error",
			err:  serrors.NewStackTraceError(errors.New("test"), stackTrace),
			want: errors.New("test").Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stackTraceError_Unwrap(t *testing.T) {
	stackTrace := serrors.StackTrace{
		{Function: "f1", File: "file1.go", Line: 1},
		{Function: "f2", File: "file2.go", Line: 2},
	}

	tests := []struct {
		name string
		err  error
		want error
	}{
		{
			name: "wrap nil error",
			err:  serrors.NewStackTraceError(nil, stackTrace),
			want: nil,
		},
		{
			name: "wrap non-nil error",
			err:  serrors.NewStackTraceError(errors.New("test"), stackTrace),
			want: errors.New("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errors.Unwrap(tt.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unwrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_withStackTrace(t *testing.T) {
	type want uint
	const (
		wantNil     want = 1
		wantWrapped want = 2
		wantSame    want = 3
	)

	tests := []struct {
		name string
		err  error
		skip int
		want want
	}{
		{
			name: "nil",
			err:  nil,
			want: wantNil,
		},
		{
			name: "want wrapped err",
			err:  errors.New("test"),
			want: wantWrapped,
		},
		{
			name: "want same",
			err:  serrors.NewStackTraceError(errors.New("test"), serrors.GetCurrentStackTrace()),
			want: wantSame,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := serrors.WithStackTrace(tt.err, tt.skip)
			switch tt.want {
			case wantNil:
				if got != nil {
					t.Errorf("want nil, got %v", got)
				}
			case wantWrapped:
				if !errors.Is(got, tt.err) {
					t.Errorf("want %v, got %v", tt.err, got)
				}
			case wantSame:
				if !reflect.DeepEqual(got, tt.err) {
					t.Errorf("want %v, got %v", tt.err, got)
				}
			default:
				t.Errorf("invalid test case")
			}
		})
	}
}

func Test_getStackTraceError(t *testing.T) {
	stackTrace := serrors.StackTrace{
		{Function: "f1", File: "file1.go", Line: 1},
		{Function: "f2", File: "file2.go", Line: 2},
	}

	wrapped := serrors.NewStackTraceError(errors.New("wrapped"), stackTrace)

	tests := []struct {
		name        string
		err         error
		wantSame    bool
		wantWrapped bool
		wantNil     bool
	}{
		{
			name:     "get same (New)",
			err:      serrors.NewStackTraceError(errors.New("test"), stackTrace),
			wantSame: true,
		},
		{
			name:     "get same (WithStackTrace)",
			err:      serrors.WithStackTrace(errors.New("test"), 1),
			wantSame: true,
		},
		{
			name:        "get wrapped (WithStackTrace with wrapped)",
			err:         serrors.WithStackTrace(wrapped, 1),
			wantWrapped: true,
		},
		{
			name:        "get wrapped (WithStackTrace with fmt.Errorf with wrapped)",
			err:         serrors.WithStackTrace(fmt.Errorf("%w", wrapped), 1),
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
			got := serrors.GetStackTraceError(tt.err)

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
	stackTrace := serrors.StackTrace{
		{Function: "f1", File: "file1.go", Line: 1},
		{Function: "f2", File: "file2.go", Line: 2},
	}

	tests := []struct {
		name        string
		err         error
		want        serrors.StackTrace
		wantCreated bool
	}{
		{
			name: "expect nil",
			err:  nil,
			want: nil,
		},
		{
			name: "expect stackTrace",
			err:  serrors.NewStackTraceError(errors.New("test"), stackTrace),
			want: stackTrace,
		},
		{
			name: "expect stackTrace (wrapped)",
			err:  fmt.Errorf("%w", serrors.NewStackTraceError(errors.New("test"), stackTrace)),
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
				if got := serrors.GetStackTrace(tt.err); got == nil {
					t.Errorf("GetStackTrace() = %v, want created stack trace", got)
				}
			} else if got := serrors.GetStackTrace(tt.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStackTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAttachedStackTrace(t *testing.T) {
	stackTrace := serrors.StackTrace{
		{Function: "f1", File: "file1.go", Line: 1},
		{Function: "f2", File: "file2.go", Line: 2},
	}

	tests := []struct {
		name string
		err  error
		want serrors.StackTrace
	}{
		{
			name: "expect nil and false (err is nil)",
			err:  nil,
			want: nil,
		},
		{
			name: "expect stackTrace and true",
			err:  serrors.NewStackTraceError(errors.New("test"), stackTrace),
			want: stackTrace,
		},
		{
			name: "expect stackTrace (wrapped) and true",
			err:  fmt.Errorf("%w", serrors.NewStackTraceError(errors.New("test"), stackTrace)),
			want: stackTrace,
		},
		{
			name: "expect nil and false (err is not stackTraceError)",
			err:  errors.New("test"),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := serrors.GetAttachedStackTrace(tt.err)

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

func TestGetStackTraces(t *testing.T) {
	newStackTrace := func(count int) serrors.StackTrace {
		ret := make([]serrors.Frame, count)
		for i := 0; i < count; i++ {
			ret[i] = serrors.Frame{Function: fmt.Sprintf("f%d", i), File: fmt.Sprintf("file%d.go", i), Line: i + 1}
		}
		return ret
	}

	serr1 := serrors.NewStackTraceError(errors.New("test1"), newStackTrace(1))
	serr2 := serrors.NewStackTraceError(errors.New("test2"), newStackTrace(2))
	serr3 := serrors.NewStackTraceError(errors.New("test3"), newStackTrace(3))

	tests := []struct {
		name            string
		err             error
		wantErrs        []error
		wantStackTraces []serrors.StackTrace
	}{
		{
			name:            "nil",
			err:             nil,
			wantErrs:        nil,
			wantStackTraces: nil,
		},
		{
			name:            "does not have stack trace for single error",
			err:             errors.New("test"),
			wantErrs:        nil,
			wantStackTraces: nil,
		},
		{
			name:            "does not have stack trace for multiple errors",
			err:             fmt.Errorf("%w", fmt.Errorf("%w", errors.New("test"))),
			wantErrs:        nil,
			wantStackTraces: nil,
		},
		{
			name:            "wrap nil error",
			err:             fmt.Errorf("%w", error(nil)),
			wantErrs:        []error{},
			wantStackTraces: nil,
		},
		{
			name:            "single stack trace",
			err:             serr1,
			wantErrs:        []error{errors.New(serr1.Error())},
			wantStackTraces: []serrors.StackTrace{newStackTrace(1)},
		},
		{
			name: "multiple stack traces",
			err:  errors.Join(serr1, serr2, serr3),
			wantErrs: []error{
				errors.New(serr1.Error()),
				errors.New(serr2.Error()),
				errors.New(serr3.Error()),
			},
			wantStackTraces: []serrors.StackTrace{
				newStackTrace(1),
				newStackTrace(2),
				newStackTrace(3),
			},
		},
		{
			name: "wrap joined errors by fmt.Errorf",
			err:  fmt.Errorf("%w", errors.Join(serr1, serr2, serr3)),
			wantErrs: []error{
				errors.New(serr1.Error()),
				errors.New(serr2.Error()),
				errors.New(serr3.Error()),
			},
			wantStackTraces: []serrors.StackTrace{
				newStackTrace(1),
				newStackTrace(2),
				newStackTrace(3),
			},
		},
		{
			name: "1 error does not have stack trace",
			err:  errors.Join(serr1, errors.New("test2"), serr3),
			wantErrs: []error{
				errors.New(serr1.Error()),
				errors.New(serr3.Error()),
			},
			wantStackTraces: []serrors.StackTrace{
				newStackTrace(1),
				newStackTrace(3),
			},
		},
		{
			name: "contains nil error in multiple errors",
			err:  &multiWrapError{errs: []error{serr1, nil, serr3}},
			wantErrs: []error{
				errors.New(serr1.Error()),
				errors.New(serr3.Error()),
			},
			wantStackTraces: []serrors.StackTrace{
				newStackTrace(1),
				newStackTrace(3),
			},
		},
		{
			name: "nested stack traces",
			err:  serrors.NewStackTraceError(serr1, newStackTrace(2)),
			wantErrs: []error{
				serr1,
			},
			wantStackTraces: []serrors.StackTrace{
				newStackTrace(2),
				// do not unwrap serr1
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.wantErrs) != len(tt.wantStackTraces) {
				t.Fatalf("want length of wantErrs and wantStackTraces is not equal")
			}

			idx := -1
			for err, stackTrace := range serrors.GetStackTraces(tt.err) {
				idx++
				if len(tt.wantErrs) <= idx || len(tt.wantStackTraces) <= idx {
					t.Fatalf("unexpected error %d: %v", idx, err)
				}

				if !reflect.DeepEqual(err, tt.wantErrs[idx]) {
					t.Errorf("error at index %d mismatch: got %v, want %v", idx, err, tt.wantErrs[idx])
				}

				if !reflect.DeepEqual(stackTrace, tt.wantStackTraces[idx]) {
					t.Errorf("unexpected stack trace %d: %v", idx, stackTrace)
				}
			}

			switch {
			case idx == -1 && 0 < len(tt.wantErrs):
				t.Errorf("no errors returned")
			case idx != -1 && idx < len(tt.wantErrs)-1:
				t.Errorf("some errors returned: %v", tt.wantErrs[idx:])
			}
		})
	}

	t.Run("break iterator", func(t *testing.T) {
		err := errors.Join(serr1, serr2, serr3)
		idx := -1

		for err, stackTrace := range serrors.GetStackTraces(err) {
			idx++
			if idx == 0 {
				if !reflect.DeepEqual(err, errors.New(serr1.Error())) {
					t.Errorf("unexpected error %d: %v", idx, err)
				}
				if !reflect.DeepEqual(stackTrace, newStackTrace(1)) {
					t.Errorf("unexpected stack trace %d: %v", idx, stackTrace)
				}
			} else {
				t.Errorf("unexpected iteration %d: %v", idx, err)
			}
			break
		}
	})
}
