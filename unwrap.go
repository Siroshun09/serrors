package serrors

const maxUnwrapCount = 256

func unwrapAll[E interface{ Unwrap() error }](err error, unwrapSelf bool, yield func(E) bool, count *int) bool {
	if count == nil {
		count = new(int)
	}

	*count++
	if maxUnwrapCount < *count {
		return false
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
		return unwrapAll[E](u, unwrapSelf, yield, count)
	case interface{ Unwrap() error }:
		u := e.Unwrap()
		if u == nil {
			return true
		}
		return unwrapAll[E](u, unwrapSelf, yield, count)
	case interface{ Unwrap() []error }:
		for _, u := range e.Unwrap() {
			if u == nil {
				continue
			}
			if !unwrapAll[E](u, unwrapSelf, yield, count) {
				return false
			}
		}
		return true
	default:
		return true
	}
}
