package serrors_test

import (
	"errors"
	"log/slog"
	"reflect"
	"testing"

	"github.com/Siroshun09/serrors/v2"
)

func TestNew(t *testing.T) {
	msg := "test error"
	attr1 := slog.String("k1", "v1")
	attr2 := slog.Int("k2", 2)

	t.Run("without attrs", func(t *testing.T) {
		err := serrors.New(msg)
		if err.Error() != msg {
			t.Errorf("Error() = %q, want %q", err.Error(), msg)
		}

		st, ok := serrors.GetAttachedStackTrace(err)
		if !ok || len(st) == 0 {
			t.Errorf("GetAttachedStackTrace() = (%v, %v), want not empty and true", st, ok)
		}

		if st[0].Function != "github.com/Siroshun09/serrors/v2_test.TestNew.func1" {
			t.Errorf("stack trace does not contain current function: %v", st)
		}

		for _, attr := range serrors.GetAttrs(err) {
			t.Errorf("unexpected attribute: %v", attr)
		}
	})

	t.Run("with attrs", func(t *testing.T) {
		err := serrors.New(msg, attr1, attr2)
		if err.Error() != msg {
			t.Errorf("Error() = %q, want %q", err.Error(), msg)
		}

		st, ok := serrors.GetAttachedStackTrace(err)
		if !ok || len(st) == 0 {
			t.Errorf("GetAttachedStackTrace() = (%v, %v), want not empty and true", st, ok)
		}

		var gotAttrs []slog.Attr
		for _, attr := range serrors.GetAttrs(err) {
			gotAttrs = append(gotAttrs, attr)
		}

		wantAttrs := []slog.Attr{attr1, attr2}
		if !reflect.DeepEqual(gotAttrs, wantAttrs) {
			t.Errorf("GetAttrs() = %v, want %v", gotAttrs, wantAttrs)
		}
	})
}

func TestWrap(t *testing.T) {
	baseErr := errors.New("base error")
	attr1 := slog.String("k1", "v1")

	t.Run("nil error", func(t *testing.T) {
		if got := serrors.Wrap(nil); got != nil {
			t.Errorf("Wrap(nil) = %v, want nil", got)
		}
	})

	t.Run("without attrs", func(t *testing.T) {
		err := serrors.Wrap(baseErr)
		if !errors.Is(err, baseErr) {
			t.Errorf("errors.Is(err, baseErr) = false, want true")
		}

		st, ok := serrors.GetAttachedStackTrace(err)
		if !ok || len(st) == 0 {
			t.Errorf("GetAttachedStackTrace() = (%v, %v), want not empty and true", st, ok)
		}

		if st[0].Function != "github.com/Siroshun09/serrors/v2_test.TestWrap.func2" {
			t.Errorf("stack trace does not contain current function: %v", st)
		}

		for _, attr := range serrors.GetAttrs(err) {
			t.Errorf("unexpected attribute: %v", attr)
		}
	})

	t.Run("with attrs", func(t *testing.T) {
		err := serrors.Wrap(baseErr, attr1)
		if !errors.Is(err, baseErr) {
			t.Errorf("errors.Is(err, baseErr) = false, want true")
		}

		st, ok := serrors.GetAttachedStackTrace(err)
		if !ok || len(st) == 0 {
			t.Errorf("GetAttachedStackTrace() = (%v, %v), want not empty and true", st, ok)
		}

		var gotAttrs []slog.Attr
		for _, attr := range serrors.GetAttrs(err) {
			gotAttrs = append(gotAttrs, attr)
		}

		wantAttrs := []slog.Attr{attr1}
		if !reflect.DeepEqual(gotAttrs, wantAttrs) {
			t.Errorf("GetAttrs() = %v, want %v", gotAttrs, wantAttrs)
		}
	})

	t.Run("wrap already has stack trace", func(t *testing.T) {
		err1 := serrors.New("err1")
		st1, _ := serrors.GetAttachedStackTrace(err1)

		err2 := serrors.Wrap(err1)
		st2, _ := serrors.GetAttachedStackTrace(err2)

		if !reflect.DeepEqual(st1, st2) {
			t.Errorf("stack trace should be reused")
		}
	})
}
