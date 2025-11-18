package serrors

import (
	"errors"
	"fmt"
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

var unknownFuncInfo = FuncInfo{
	Name: "UNKNOWN",
}

func newStackTraceFromCallers(skip int) StackTrace {
	pcs := make([]uintptr, 64)
	l := runtime.Callers(skip+2, pcs) // callers -> newStackTraceFromCallers

	stacks := make(StackTrace, l)
	for i, pc := range pcs[:l] {
		stacks[i] = newFuncInfo(pc, runtime.FuncForPC(pc))
	}

	return stacks
}

func newFuncInfo(pc uintptr, f *runtime.Func) FuncInfo {
	if f == nil {
		return unknownFuncInfo
	} else {
		file, line := f.FileLine(pc)
		return FuncInfo{
			Name: f.Name(),
			File: file,
			Line: line,
		}
	}
}
