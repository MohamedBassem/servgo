package servgo

type UnparsableRequestError struct {
	s string
}

func (err *UnparsableRequestError) Error() string {
	return err.s
}
