package serrors

func NewStackTrace(skip int, limit int) StackTrace {
	return newStackTrace(skip, limit)
}

func NewStackTraceError(err error, st StackTrace) error {
	return &stackTraceError{err: err, st: st}
}

func WithStackTrace(err error, skip int) error {
	return withStackTrace(err, skip)
}

type StackTraceError = *stackTraceError

func GetStackTraceError(err error) StackTraceError {
	return getStackTraceError(err)
}
