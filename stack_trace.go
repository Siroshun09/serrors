package serrors

import (
	"errors"
	"iter"
	"runtime"
	"strconv"
)

type Frame struct {
	Function string
	File     string
	Line     int
}

// String formats Frame as "function (file:line)"
func (s Frame) String() string {
	ret, _ := s.AppendText(nil) // AppendText never returns error
	return string(ret)
}

func (s Frame) AppendText(ret []byte) ([]byte, error) {
	ret = append(ret, s.Function...)
	ret = append(ret, ' ', '(')
	ret = append(ret, s.File...)
	ret = append(ret, ':')
	ret = strconv.AppendInt(ret, int64(s.Line), 10)
	ret = append(ret, ')')
	return ret, nil
}

type StackTrace []Frame

func (st StackTrace) String() string {
	buf := make([]byte, 0, 1024)
	buf, _ = st.AppendText(buf) // AppendText never returns error
	return string(buf)
}

func (st StackTrace) AppendText(ret []byte) ([]byte, error) {
	for i, frame := range st {
		if i != 0 {
			ret = append(ret, '\n')
		}
		ret, _ = frame.AppendText(ret) // AppendText never returns error
	}
	return ret, nil
}

const defaultStackTraceLimit = 64

// GetCurrentStackTrace returns the current StackTrace.
func GetCurrentStackTrace() StackTrace {
	return newStackTrace(1, defaultStackTraceLimit)
}

func newStackTrace(skip int, limit int) StackTrace {
	pcs := make([]uintptr, limit)
	l := runtime.Callers(skip+2, pcs) // Callers -> newStackTrace
	frames := runtime.CallersFrames(pcs[:l])
	st := make(StackTrace, 0, l)

	for {
		frame, more := frames.Next()
		st = append(st, Frame{
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		})
		if !more {
			break
		}
	}

	return st
}

type stackTraceError struct {
	err error
	st  StackTrace
}

func (e *stackTraceError) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e *stackTraceError) Unwrap() error {
	return e.err
}

func withStackTrace(err error, skip int) error {
	if err == nil {
		return nil
	}

	serr := getStackTraceError(err)
	if serr != nil {
		return err
	}

	return &stackTraceError{
		err: err,
		st:  newStackTrace(skip+1, defaultStackTraceLimit),
	}
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

	return newStackTrace(1, defaultStackTraceLimit) // GetStackTrace
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
		return serr.st, true
	}

	return nil, false
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
		unwrapAll[*stackTraceError](err, false, func(sterr *stackTraceError) bool {
			return yield(sterr.err, sterr.st)
		}, new(int))
	}
}
