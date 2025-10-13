package todo

import (
	"context"
)

// Repository defines the contract for todo persistence operations
type Repository interface {
	// Create saves a new todo to the repository
	Create(ctx context.Context, todo *Todo) error

	// GetByID retrieves a todo by its ID
	GetByID(ctx context.Context, id TodoID) (*Todo, error)

	// GetAll retrieves all todos
	GetAll(ctx context.Context) ([]*Todo, error)

	// Update updates an existing todo
	Update(ctx context.Context, todo *Todo) error

	// Delete removes a todo by its ID
	Delete(ctx context.Context, id TodoID) error

	// FindByCompleted retrieves todos filtered by completion status
	FindByCompleted(ctx context.Context, completed bool) ([]*Todo, error)

	// Exists checks if a todo exists by its ID
	Exists(ctx context.Context, id TodoID) (bool, error)
}

// RepositoryError represents errors that can occur in repository operations
type RepositoryError struct {
	Operation string
	Err       error
}

func (e *RepositoryError) Error() string {
	return e.Operation + ": " + e.Err.Error()
}

func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// NewRepositoryError creates a new repository error
func NewRepositoryError(operation string, err error) *RepositoryError {
	return &RepositoryError{
		Operation: operation,
		Err:       err,
	}
}
