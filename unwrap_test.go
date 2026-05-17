package serrors_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/Siroshun09/serrors/v2"
)

type mockE struct {
	msg string
	err error
}

func (e *mockE) Error() string {
	if e.err != nil {
		return e.msg + ": " + e.err.Error()
	}
	return e.msg
}

func (e *mockE) Unwrap() error {
	return e.err
}

func TestUnwrapAll(t *testing.T) {
	e1 := &mockE{msg: "e1"}
	e2 := &mockE{msg: "e2"}
	e3 := &mockE{msg: "e3"}

	tests := []struct {
		name       string
		err        error
		unwrapSelf bool
		yield      func(*mockE) bool
		countRef   *int
		want       bool
		wantErrs   []*mockE
		wantCount  int
	}{
		{
			name:       "nil error",
			err:        nil,
			unwrapSelf: false,
			yield: func(e *mockE) bool {
				return true
			},
			countRef:  new(int),
			want:      true,
			wantErrs:  nil,
			wantCount: 1,
		},
		{
			name:       "not E and no Unwrap",
			err:        errors.New("base"),
			unwrapSelf: false,
			yield: func(e *mockE) bool {
				return true
			},
			countRef:  new(int),
			want:      true,
			wantErrs:  nil,
			wantCount: 1,
		},
		{
			name:       "is E, yield returns false",
			err:        e1,
			unwrapSelf: false,
			yield: func(e *mockE) bool {
				return false
			},
			countRef:  new(int),
			want:      false,
			wantErrs:  []*mockE{e1},
			wantCount: 1,
		},
		{
			name:       "is E, count is nil",
			err:        e1,
			unwrapSelf: false,
			yield: func(e *mockE) bool {
				return false
			},
			countRef:  nil,
			want:      false,
			wantErrs:  []*mockE{e1},
			wantCount: 1,
		},
		{
			name:       "is E, unwrapSelf is false",
			err:        &mockE{msg: "e1", err: e2},
			unwrapSelf: false,
			yield: func(e *mockE) bool {
				return true
			},
			countRef:  new(int),
			want:      true,
			wantErrs:  []*mockE{{msg: "e1", err: e2}},
			wantCount: 1,
		},
		{
			name:       "is E, unwrapSelf is true, Unwrap returns nil",
			err:        e1,
			unwrapSelf: true,
			yield: func(e *mockE) bool {
				return true
			},
			countRef:  new(int),
			want:      true,
			wantErrs:  []*mockE{e1},
			wantCount: 1,
		},
		{
			name:       "is E, unwrapSelf is true, recursion",
			err:        &mockE{msg: "e1", err: e2},
			unwrapSelf: true,
			yield: func(e *mockE) bool {
				return true
			},
			countRef:  new(int),
			want:      true,
			wantErrs:  []*mockE{{msg: "e1", err: e2}, e2},
			wantCount: 2,
		},
		{
			name:       "not E, has Unwrap() error, Unwrap returns nil",
			err:        fmt.Errorf("wrap: %w", error(nil)),
			unwrapSelf: false,
			yield: func(e *mockE) bool {
				return true
			},
			countRef:  new(int),
			want:      true,
			wantErrs:  nil,
			wantCount: 1,
		},
		{
			name:       "not E, has Unwrap() error, recursion finds E",
			err:        fmt.Errorf("wrap: %w", e1),
			unwrapSelf: false,
			yield: func(e *mockE) bool {
				return true
			},
			countRef:  new(int),
			want:      true,
			wantErrs:  []*mockE{e1},
			wantCount: 2,
		},
		{
			name:       "has Unwrap() []error, contains nil",
			err:        &multiWrapError{errs: []error{nil, e1, nil, e2}},
			unwrapSelf: false,
			yield: func(e *mockE) bool {
				return true
			},
			countRef:  new(int),
			want:      true,
			wantErrs:  []*mockE{e1, e2},
			wantCount: 3,
		},
		{
			name:       "has Unwrap() []error, recursion stops",
			err:        errors.Join(e1, e2),
			unwrapSelf: false,
			yield: func(e *mockE) bool {
				return false
			},
			countRef:  new(int),
			want:      false,
			wantErrs:  []*mockE{e1},
			wantCount: 2,
		},
		{
			name:       "has Unwrap() []error, all processed",
			err:        errors.Join(e1, e2),
			unwrapSelf: false,
			yield: func(e *mockE) bool {
				return true
			},
			countRef:  new(int),
			want:      true,
			wantErrs:  []*mockE{e1, e2},
			wantCount: 3,
		},
		{
			name:       "has Unwrap() []error, count equal to maxUnwrapCount",
			err:        errors.Join(e1, e2, e3),
			unwrapSelf: false,
			yield: func(e *mockE) bool {
				return true
			},
			countRef: func() *int {
				initialCount := 253 // 254:self -> 255:e1 -> 256:e2
				return &initialCount
			}(),
			want:      false,
			wantErrs:  []*mockE{e1, e2},
			wantCount: 257,
		},
		{
			name:       "has Unwrap() []error, count exceeds maxUnwrapCount (start 254)",
			err:        errors.Join(e1, e2),
			unwrapSelf: false,
			yield: func(e *mockE) bool {
				return true
			},
			countRef: func() *int {
				initialCount := 254 // 255:self -> 256:e1 -> (stopped) 257:e2
				return &initialCount
			}(),
			want:      false,
			wantErrs:  []*mockE{e1},
			wantCount: 257,
		},
		{
			name:       "has Unwrap() []error, count exceeds maxUnwrapCount (start 255)",
			err:        errors.Join(e1, e2),
			unwrapSelf: false,
			yield: func(e *mockE) bool {
				return true
			},
			countRef: func() *int {
				initialCount := 255 // 256:self -> (stopped) 257:e1 -> 258:e2
				return &initialCount
			}(),
			want:      false,
			wantErrs:  nil,
			wantCount: 257,
		},
		{
			name:       "has Unwrap() []error, count exceeds maxUnwrapCount (start 256)",
			err:        errors.Join(e1, e2),
			unwrapSelf: false,
			yield: func(e *mockE) bool {
				return true
			},
			countRef: func() *int {
				initialCount := 256 // (stopped) 257:self -> 258:e1 -> 259:e2
				return &initialCount
			}(),
			want:      false,
			wantErrs:  nil,
			wantCount: 257,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotErrs []*mockE
			yield := func(e *mockE) bool {
				gotErrs = append(gotErrs, e)
				return tt.yield(e)
			}

			got := serrors.UnwrapAll[*mockE](tt.err, tt.unwrapSelf, yield, tt.countRef)
			if got != tt.want {
				t.Errorf("unwrapAll() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("yielded errors = %v, want %v", gotErrs, tt.wantErrs)
			}

			if tt.countRef != nil && *tt.countRef != tt.wantCount {
				t.Errorf("count = %d, want %d", *tt.countRef, tt.wantCount)
			}
		})
	}

	t.Run("too many wrapped error stops at maxUnwrapCount", func(t *testing.T) {
		var err error = &mockE{msg: "base"}
		for i := 0; i < 300; i++ {
			err = &mockE{msg: fmt.Sprintf("level%d", i), err: err}
		}

		var count int
		yield := func(e *mockE) bool {
			count++
			return true
		}

		got := serrors.UnwrapAll[*mockE](err, true, yield, new(int))
		if got {
			t.Errorf("unwrapAll() = %v, want false", got)
		}

		if count != 256 {
			t.Errorf("yield called %d times, want 256", count)
		}
	})
}
