package usecase

import "errors"

type UseCaseError struct { //nolint revive
	err error
}

func (e UseCaseError) Error() string {
	return e.err.Error()
}

func (e UseCaseError) Unwrap() error {
	return e.err
}

func NewUseCaseErrorFromMsg(msg string) UseCaseError {
	return UseCaseError{err: errors.New(msg)}
}

func NewUseCaseError(err error) UseCaseError {
	return UseCaseError{err: err}
}

func IsUseCaseError(err error) bool {
	if err == nil {
		return false
	}
	return errors.As(err, &UseCaseError{})
}

var ErrNotFound = errors.New("not found")
