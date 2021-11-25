package utils

type ErrorGroup []error

func (e *ErrorGroup) Add(errs ...error) {
	*e = append(*e, errs...)
}

func (e *ErrorGroup) Count() int {
	return len(*e)
}

func (e *ErrorGroup) List() []error {
	return *e
}
