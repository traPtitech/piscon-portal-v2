package usecase

import "errors"

type ErrBadRequest struct {
	err error
}

func (e ErrBadRequest) Error() string {
	return e.err.Error()
}

func NewErrBadRequest(msg string) ErrBadRequest {
	return ErrBadRequest{err: errors.New(msg)}
}

func NewErrBadRequestFromErr(err error) ErrBadRequest {
	return ErrBadRequest{err: err}
}

func IsErrBadRequest(err error) bool {
	if err == nil {
		return false
	}
	return errors.As(err, &ErrBadRequest{})
}

var ErrNotFound = errors.New("not found")
