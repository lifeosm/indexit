package exitcode

import "errors"

const (
	OK    = 0
	Fail  = 1
	Usage = 2
)

type Error struct {
	Code int
	Err  error
}

func New(code int, err error) *Error {
	return &Error{Code: code, Err: err}
}

func (e *Error) Error() string {
	if e == nil || e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func FromError(err error) int {
	if err == nil {
		return OK
	}

	var exit *Error
	if errors.As(err, &exit) && exit.Code != OK {
		return exit.Code
	}
	return Fail
}
