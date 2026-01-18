package errors

import "fmt"

type ApplicationError struct {
	Code       string
	Message    string
	StatusCode int
	Err        error
}

func (e *ApplicationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Allows error.Is() to work
func (e *ApplicationError) Unwrap() error {
	return e.Err
}

func NewInvalidInputError(err error) *ApplicationError {
	return &ApplicationError{
		Code:       "INVALID_INPUT",
		Message:    "invalid request data",
		StatusCode: 400,
		Err:        err,
	}
}

func NewDuplicateCPFError() *ApplicationError {
	return &ApplicationError{
		Code:       "DUPLICATE_CPF",
		Message:    "CPF already registered",
		StatusCode: 409,
	}
}

func NewInternalError(message string, err error) *ApplicationError {
	return &ApplicationError{
		Code:       "INTERNAL_ERROR",
		Message:    message,
		StatusCode: 500,
		Err:        err,
	}
}
