package serrors_test

type multiWrapError struct {
	errs []error // possibly contains nil
}

func (e *multiWrapError) Error() string {
	return "multi"
}

func (e *multiWrapError) Unwrap() []error {
	return e.errs
}
