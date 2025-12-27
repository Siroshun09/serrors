package serrors

import (
	"errors"
	"fmt"
	"iter"
	"runtime"
	"strconv"
	"strings"
)

type stackTraceError struct {
	err        error
	stackTrace StackTrace
}

func (e *stackTraceError) Error() string {
	return e.err.Error()
}

func (e *stackTraceError) Unwrap() error {
	return e.err
}

// New creates an error with a StackTrace.
func New(msg string) error {
	return withStackTrace(errors.New(msg))
}

// Errorf creates an error with a StackTrace.
func Errorf(format string, args ...any) error {
	return withStackTrace(fmt.Errorf(format, args...))
}

// WithStackTrace creates an error with a StackTrace.
//
// If err already has a StackTrace, this function returns err as-is.
//
// Also, if err is nil, this function returns nil.
func WithStackTrace(err error) error {
	return withStackTrace(err)
}

func withStackTrace(err error) error {
	if err == nil {
		return nil
	}

	serr := getStackTraceError(err)
	if serr != nil {
		return err
	}

	return &stackTraceError{
		err:        err,
		stackTrace: newStackTraceFromCallers(2), // withStackTrace -> caller (New/Errorf/WithStackTrace)
	}
}

// GetStackTrace returns a StackTrace for err.
//
// If err does not have a StackTrace, this function creates the current StackTrace.
//
// Also, if err is nil, this function returns nil.
func GetStackTrace(err error) StackTrace {
	if err == nil {
		return nil
	}

	stackTrace, ok := GetAttachedStackTrace(err)
	if ok {
		return stackTrace
	}

	return newStackTraceFromCallers(1) // GetStackTrace
}

// GetAttachedStackTrace returns the StackTrace if the given error has one.
//
// The returned bool indicates whether the given error has a StackTrace.
func GetAttachedStackTrace(err error) (StackTrace, bool) {
	if err == nil {
		return nil, false
	}

	serr := getStackTraceError(err)
	if serr != nil {
		return serr.stackTrace, true
	}

	return nil, false
}

// GetCurrentStackTrace returns the current StackTrace.
func GetCurrentStackTrace() StackTrace {
	return newStackTraceFromCallers(1)
}

func getStackTraceError(err error) *stackTraceError {
	if err == nil {
		return nil
	}

	var serr *stackTraceError
	if !errors.As(err, &serr) {
		return nil
	}

	return serr
}

// FuncInfo contains values obtained from runtime.Func.
type FuncInfo struct {
	// Name is the name of the function.
	Name string
	// File is the source file where the function is defined.
	File string
	// Line is the line number in that file.
	Line int
}

// String formats FuncInfo as "name (file:line)"
func (s FuncInfo) String() string {
	return s.Name + " (" + s.File + ":" + strconv.Itoa(s.Line) + ")"
}

// StackTrace is an array of FuncInfo.
type StackTrace []FuncInfo

// String formats StackTrace using FuncInfo.String.
func (s StackTrace) String() string {
	builder := strings.Builder{}
	for i, stack := range s {
		if 0 < i {
			builder.WriteString("\n")
		}
		builder.WriteString(stack.String())
	}
	return builder.String()
}

func newStackTraceFromCallers(skip int) StackTrace {
	pcs := make([]uintptr, 64)
	l := runtime.Callers(skip+2, pcs) // callers -> newStackTraceFromCallers
	frames := runtime.CallersFrames(pcs[:l])
	st := make(StackTrace, 0, l)

	for {
		frame, more := frames.Next()
		st = append(st, FuncInfo{
			Name: frame.Function,
			File: frame.File,
			Line: frame.Line,
		})
		if !more {
			break
		}
	}

	return st
}

// GetStackTraces returns a sequence of errors and their associated StackTrace.
//
// This function recursively returns errors and stack traces by repeating the following process:
//
// 1. If the given error has a StackTrace itself, this function returns the wrapped error and its StackTrace
//   - In this case, the iterator returns only one error.
//
// 2. If the given error has an Unwrap() error function, this function calls it and tries to process step 1 again
// 3. If the given error has an Unwrap() []error function, this function calls it and tries to process step 1 for each error in the returned slice
func GetStackTraces(err error) iter.Seq2[error, StackTrace] {
	return func(yield func(error, StackTrace) bool) {
		tryYieldStackTrace(err, yield)
	}
}

func tryYieldStackTrace(err error, yield func(error, StackTrace) bool) bool {
	switch x := err.(type) {
	case *stackTraceError:
		return yield(x.err, x.stackTrace)
	case interface{ Unwrap() error }:
		err = x.Unwrap()
		if err == nil {
			return true
		}
		return tryYieldStackTrace(err, yield)
	case interface{ Unwrap() []error }:
		for _, err := range x.Unwrap() {
			if err == nil {
				continue
			}
			if !tryYieldStackTrace(err, yield) {
				return false
			}
		}
		return true
	default:
		return true
	}
}
