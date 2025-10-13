package todo

import "fmt"

// UseCaseError represents errors that can occur in use case operations
type UseCaseError struct {
	Operation string
	Err       error
}

func (e *UseCaseError) Error() string {
	return fmt.Sprintf("%s: %v", e.Operation, e.Err)
}

func (e *UseCaseError) Unwrap() error {
	return e.Err
}

// NewUseCaseError creates a new use case error
func NewUseCaseError(operation string, err error) *UseCaseError {
	return &UseCaseError{
		Operation: operation,
		Err:       err,
	}
}
