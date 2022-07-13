package domain

import "fmt"

type ErrorType int

const (
	InvalidInput ErrorType = iota
	ResourceNotFound
	SystemError
	ResourceConflict
	NotImplemented

	errInvalidInputDataStr = "invalid input data"
	errResourceNotFoundStr = "resource not found"
	errSystemErrorStr      = "unexpected system error"
	errResourceConflictStr = "resource state conflict"
	errNotImplementedStr   = "not implemented"
)

var domainErrors = map[ErrorType]string{
	InvalidInput:     errInvalidInputDataStr,
	ResourceNotFound: errResourceNotFoundStr,
	SystemError:      errSystemErrorStr,
	ResourceConflict: errResourceConflictStr,
	NotImplemented:   errNotImplementedStr,
}

func (errT ErrorType) String() string {
	errStr, ok := domainErrors[errT]
	if !ok {
		return fmt.Sprintf("unknown error type [%d]", errT)
	}

	return errStr
}

// Error is a custom error type used to represent errors that occur in the domain
type Error struct {
	Code ErrorType
	Msg  string
	Err  error
}

// newInvalidInputError is a helper function that constructs a new domainError of code invalidInputData with the given message and error.
func newInvalidInputError(msg string, err error) *Error {
	return &Error{Code: InvalidInput, Msg: msg, Err: err}
}

// newResourceNotFoundError is a helper function that constructs a new domainError of code resourceNotFound with the given message and error.
func newResourceNotFoundError(msg string, err error) *Error {
	return &Error{Code: ResourceNotFound, Msg: msg, Err: err}
}

// newUnexpectedSystemError is a helper function that constructs a new domainError of code unexpectedSystemError with the given message and error.
func newSystemError(msg string, err error) *Error {
	return &Error{Code: SystemError, Msg: msg, Err: err}
}

// newResourceConflictError is a helper function that constructs a new domainError of code resourceConflict with the given message and error.
func newResourceConflictError(msg string, err error) *Error {
	return &Error{Code: ResourceConflict, Msg: msg, Err: err}
}

// newNotImplementedError is a helper function that constructs a new domainError of code notImplemented with the given message and error.
func newNotImplementedError(msg string, err error) *Error {
	return &Error{Code: NotImplemented, Msg: msg, Err: err}
}

func (cerr *Error) Error() string {
	errMsg := fmt.Sprintf("%s - %v", domainErrors[cerr.Code], cerr.Msg)
	if cerr.Err != nil {
		errMsg += fmt.Sprintf(": [%v]", cerr.Err.Error())
	}

	return errMsg
}

// WrapMessage allows the caller to wrap the error message with additional information
func (cerr *Error) WrapMessage(msg string) *Error {
	cerr.Msg = msg + ", " + cerr.Msg
	return cerr
}
