package api

type notFoundError struct {
	mes string
}

func (e notFoundError) Error() string {
	return e.mes
}

type malformedError struct {
	mes string
}

func (e malformedError) Error() string {
	return e.mes
}

type conflictError struct {
	mes string
}

func (e conflictError) Error() string {
	return e.mes
}

type unsupportedError struct {
	mes string
}

func (e unsupportedError) Error() string {
	return e.mes
}
