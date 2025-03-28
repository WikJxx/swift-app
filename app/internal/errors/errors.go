package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return e.Message
}

func New(message string, statusCode int) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: statusCode,
	}
}

var (
	ErrBadRequest = New("bad request", http.StatusBadRequest)
	ErrNotFound   = New("not found", http.StatusNotFound)
	ErrConflict   = New("conflict", http.StatusConflict)
	ErrInternal   = New("internal server error", http.StatusInternalServerError)
)

func Wrap(base *AppError, format string, args ...interface{}) *AppError {
	return &AppError{
		Message:    fmt.Sprintf(format, args...),
		StatusCode: base.StatusCode,
	}
}

func GetStatusCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.StatusCode
	}
	return http.StatusInternalServerError
}
