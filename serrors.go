package serrors

import (
	"errors"
	"log/slog"
)

// New creates an error with a StackTrace.
func New(msg string, attrs ...slog.Attr) error {
	sterr := withStackTrace(errors.New(msg), 1)
	if len(attrs) == 0 {
		return sterr
	}
	return withAttrs(sterr, attrs)
}

func Wrap(err error, attrs ...slog.Attr) error {
	if err == nil {
		return nil
	}

	sterr := withStackTrace(err, 1)
	if len(attrs) == 0 {
		return sterr
	}
	return withAttrs(sterr, attrs)
}
