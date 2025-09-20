package errors

import "fmt"

type ErrorType string

const (
	NotFound         ErrorType = "NOT_FOUND"
	Validation       ErrorType = "VALIDATION"
	Conflict         ErrorType = "CONFLICT"
	Internal         ErrorType = "INTERNAL"
	ExternalService  ErrorType = "EXTERNAL_SERVICE"
	Unauthorized     ErrorType = "UNAUTHORIZED"
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewNotFoundError(message string) error {
	return &AppError{
		Type:    NotFound,
		Message: message,
	}
}

func NewValidationError(message string) error {
	return &AppError{
		Type:    Validation,
		Message: message,
	}
}

func NewConflictError(message string) error {
	return &AppError{
		Type:    Conflict,
		Message: message,
	}
}

func NewInternalError(message string) error {
	return &AppError{
		Type:    Internal,
		Message: message,
	}
}

func NewExternalServiceError(message string) error {
	return &AppError{
		Type:    ExternalService,
		Message: message,
	}
}

func NewUnauthorizedError(message string) error {
	return &AppError{
		Type:    Unauthorized,
		Message: message,
	}
}

func IsNotFound(err error) bool {
	appErr, ok := err.(*AppError)
	return ok && appErr.Type == NotFound
}

func IsValidation(err error) bool {
	appErr, ok := err.(*AppError)
	return ok && appErr.Type == Validation
}

func IsConflict(err error) bool {
	appErr, ok := err.(*AppError)
	return ok && appErr.Type == Conflict
}

func IsInternal(err error) bool {
	appErr, ok := err.(*AppError)
	return ok && appErr.Type == Internal
}

func IsExternalService(err error) bool {
	appErr, ok := err.(*AppError)
	return ok && appErr.Type == ExternalService
}

func IsUnauthorized(err error) bool {
	appErr, ok := err.(*AppError)
	return ok && appErr.Type == Unauthorized
}