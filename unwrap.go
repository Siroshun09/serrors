package serrors

const maxUnwrapDepth = 256

func unwrapAll[E interface{ Unwrap() error }](err error, unwrapSelf bool, yield func(E) bool, depth *int) bool {
	if depth == nil {
		depth = new(int)
	}

	*depth++
	if maxUnwrapDepth < *depth {
		return true
	}

	switch e := err.(type) {
	case E:
		if !yield(e) {
			return false
		}

		if !unwrapSelf {
			return true
		}

		u := e.Unwrap()
		if u == nil {
			return true
		}
		return unwrapAll[E](u, unwrapSelf, yield, depth)
	case interface{ Unwrap() error }:
		u := e.Unwrap()
		if u == nil {
			return true
		}
		return unwrapAll[E](u, unwrapSelf, yield, depth)
	case interface{ Unwrap() []error }:
		for _, u := range e.Unwrap() {
			if u == nil {
				continue
			}
			if !unwrapAll[E](u, unwrapSelf, yield, depth) {
				return false
			}
		}
		return true
	default:
		return true
	}
}
