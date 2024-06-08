package errors

import (
	"fmt"
)

type AppError struct {
	code    int    // Business error code
	status  int    // HTTP status code
	message string // Error message
	stack   string // Stack information
}

func (e *AppError) Error() string {
	return e.message
}

func (e *AppError) Code() int {
	return e.code
}

func (e *AppError) Status() int {
	return e.status
}

func NewAppError(code, status int, message string, cause error) *AppError {
	return &AppError{
		code:    code,
		status:  status,
		message: fmt.Sprintf("%s: %s", message, cause.Error()),
		stack:   fmt.Sprintf("%+v", cause),
	}
}
