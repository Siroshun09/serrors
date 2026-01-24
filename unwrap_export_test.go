package serrors

func UnwrapAll[E interface{ Unwrap() error }](err error, unwrapSelf bool, yield func(E) bool) bool {
	return unwrapAll[E](err, unwrapSelf, yield, new(int))
}
