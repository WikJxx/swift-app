package errors

import "errors"

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")
	ErrInternal   = errors.New("internal server error")
)

func IsBadRequest(err error) bool {
	return errors.Is(err, ErrBadRequest)
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func IsConflict(err error) bool {
	return errors.Is(err, ErrConflict)
}

func IsInternal(err error) bool {
	return errors.Is(err, ErrInternal)
}

func MapToStatusCode(err error) int {
	switch {
	case IsBadRequest(err):
		return 400
	case IsNotFound(err):
		return 404
	case IsConflict(err):
		return 409
	case IsInternal(err):
		return 500
	default:
		return 500
	}
}
