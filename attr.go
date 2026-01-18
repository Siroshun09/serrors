package serrors

import (
	"iter"
	"log/slog"
)

type attrError struct {
	err   error
	attrs []slog.Attr
}

func (e *attrError) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e *attrError) Unwrap() error {
	return e.err
}

func withAttrs(err error, attrs []slog.Attr) error {
	return &attrError{
		err:   err,
		attrs: attrs,
	}
}

func GetAttrs(err error) iter.Seq2[error, slog.Attr] {
	return func(yield func(error, slog.Attr) bool) {
		unwrapAll[*attrError](err, true, func(attrErr *attrError) bool {
			for _, attr := range attrErr.attrs {
				if !yield(attrErr.err, attr) {
					return false
				}
			}
			return true
		}, new(int))
	}
}
