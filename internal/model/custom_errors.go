package model

import "fmt"

// CustomError defines an error with a status code and message
type CustomError struct {
	Code    int
	Message string
}

// Error implements the error interface
func (e *CustomError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// New creates a new CustomError
func NewCustomError(code int, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

var (
	ErrInternal = NewCustomError(500, "internal server error")
)
