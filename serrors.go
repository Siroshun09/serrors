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

// New creates an error that has StackTrace.
func New(msg string) error {
	return withStackTrace(errors.New(msg))
}

// Errorf creates an error that has StackTrace.
func Errorf(format string, args ...any) error {
	return withStackTrace(fmt.Errorf(format, args...))
}

// WithStackTrace creates an error that has StackTrace.
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
		return serr
	}

	return &stackTraceError{
		err:        err,
		stackTrace: newStackTraceFromCallers(2), // withStackTrace -> caller (New/Errorf/WithStackTrace)
	}
}

// GetStackTrace returns StackTrace from err.
//
// If err does not have StackTrace, this function creates current StackTrace.
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

// GetAttachedStackTrace returns StackTrace when the given error has it.
//
// The returning bool value indicates that the given error has StackTrace.
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

// FuncInfo has values that are obtained from runtime.Func.
type FuncInfo struct {
	// Name is a name of the function
	Name string
	// File is a filename of the function.
	File string
	// Line is a line number of the function.
	Line int
}

// String formats FuncInfo as "name (file:line)"
func (s FuncInfo) String() string {
	return s.Name + " (" + s.File + ":" + strconv.Itoa(s.Line) + ")"
}

// StackTrace is a array of FuncInfo.
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
