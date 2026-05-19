package serrors_test

import (
	"errors"
	"log/slog"
	"reflect"
	"testing"

	"github.com/Siroshun09/serrors/v2"
)

func Test_attrError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "wrap nil error",
			err:  serrors.WithAttrs(nil, slog.String("k", "v")),
			want: "",
		},
		{
			name: "wrap non-nil error",
			err:  serrors.WithAttrs(errors.New("test"), slog.String("k", "v")),
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

func Test_attrError_Unwrap(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want error
	}{
		{
			name: "wrap nil error",
			err:  serrors.WithAttrs(nil, slog.String("k", "v")),
			want: nil,
		},
		{
			name: "wrap non-nil error",
			err:  serrors.WithAttrs(errors.New("test"), slog.String("k", "v")),
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

func TestGetAttrs(t *testing.T) {
	attr1 := slog.String("k1", "v1")
	attr2 := slog.Int("k2", 2)
	attr3 := slog.Bool("k3", true)

	base1 := errors.New("base1")
	base2 := errors.New("base2")

	err1 := serrors.WithAttrs(base1, attr1)
	err2 := serrors.WithAttrs(base2, attr2, attr3)

	tests := []struct {
		name      string
		err       error
		wantErrs  []error
		wantAttrs []slog.Attr
	}{
		{
			name:      "nil",
			err:       nil,
			wantErrs:  nil,
			wantAttrs: nil,
		},
		{
			name:      "no attributes for single error",
			err:       errors.New("test"),
			wantErrs:  nil,
			wantAttrs: nil,
		},
		{
			name:      "single error with one attribute",
			err:       err1,
			wantErrs:  []error{base1},
			wantAttrs: []slog.Attr{attr1},
		},
		{
			name:      "single error with multiple attributes",
			err:       err2,
			wantErrs:  []error{base2, base2},
			wantAttrs: []slog.Attr{attr2, attr3},
		},
		{
			name:      "multiple errors with attributes",
			err:       errors.Join(err1, err2),
			wantErrs:  []error{base1, base2, base2},
			wantAttrs: []slog.Attr{attr1, attr2, attr3},
		},
		{
			name:      "nested attrError",
			err:       serrors.WithAttrs(err1, attr2),
			wantErrs:  []error{err1, base1},
			wantAttrs: []slog.Attr{attr2, attr1},
		},
		{
			name:      "nested attrError with multiple attributes",
			err:       serrors.WithAttrs(err2, attr1),
			wantErrs:  []error{err2, base2, base2},
			wantAttrs: []slog.Attr{attr1, attr2, attr3},
		},
		{
			name:      "wrap nil error",
			err:       serrors.WithAttrs(nil, attr1),
			wantErrs:  []error{nil},
			wantAttrs: []slog.Attr{attr1},
		},
		{
			name:      "check Error() method",
			err:       err1,
			wantErrs:  []error{base1},
			wantAttrs: []slog.Attr{attr1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.wantErrs) != len(tt.wantAttrs) {
				t.Fatalf("want length of wantErrs and wantAttrs is not equal")
			}

			idx := -1
			for err, attr := range serrors.GetAttrs(tt.err) {
				idx++
				if len(tt.wantErrs) <= idx || len(tt.wantAttrs) <= idx {
					t.Fatalf("unexpected attribute %d: %v", idx, attr)
				}

				if !reflect.DeepEqual(err, tt.wantErrs[idx]) {
					t.Errorf("error at index %d mismatch: got %v, want %v", idx, err, tt.wantErrs[idx])
				}

				if !reflect.DeepEqual(attr, tt.wantAttrs[idx]) {
					t.Errorf("attribute at index %d mismatch: got %v, want %v", idx, attr, tt.wantAttrs[idx])
				}
			}

			switch {
			case idx == -1 && 0 < len(tt.wantErrs):
				t.Errorf("no attributes returned")
			case idx != -1 && idx < len(tt.wantErrs)-1:
				t.Errorf("some attributes returned: %v", tt.wantAttrs[idx:])
			}
		})
	}

	t.Run("break iterator", func(t *testing.T) {
		err := errors.Join(err1, err2)
		idx := -1

		for err, attr := range serrors.GetAttrs(err) {
			idx++
			if idx == 0 {
				if !reflect.DeepEqual(err, base1) {
					t.Errorf("unexpected error %d: %v", idx, err)
				}
				if !reflect.DeepEqual(attr, attr1) {
					t.Errorf("unexpected attribute %d: %v", idx, attr)
				}
			} else {
				t.Errorf("unexpected iteration %d: %v", idx, err)
			}
			break
		}
	})
}
