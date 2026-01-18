package serrors

import (
	"log/slog"
)

func WithAttrs(err error, attrs ...slog.Attr) error {
	return withAttrs(err, attrs)
}
