package todo

import (
	"context"
	"errors"

	"todo-backend/internal/domain/todo"
)

// DeleteTodoRequest represents the input for deleting a todo
type DeleteTodoRequest struct {
	ID string `uri:"id" binding:"required"`
}

// DeleteTodoUseCase handles deleting existing todos
type DeleteTodoUseCase struct {
	repository     todo.Repository
	domainService  todo.Service
}

// NewDeleteTodoUseCase creates a new DeleteTodoUseCase
func NewDeleteTodoUseCase(repository todo.Repository, domainService todo.Service) *DeleteTodoUseCase {
	return &DeleteTodoUseCase{
		repository:    repository,
		domainService: domainService,
	}
}

// Execute deletes an existing todo
func (uc *DeleteTodoUseCase) Execute(ctx context.Context, request DeleteTodoRequest) error {
	// Parse todo ID
	todoID, err := todo.NewTodoIDFromString(request.ID)
	if err != nil {
		return NewUseCaseError("invalid todo ID", err)
	}

	// Check if the todo can be deleted using domain service
	if err := uc.domainService.CanDelete(ctx, todoID); err != nil {
		if errors.Is(err, errors.New("todo not found")) {
			return NewUseCaseError("todo not found", err)
		}
		return NewUseCaseError("cannot delete todo", err)
	}

	// Delete the todo
	if err := uc.repository.Delete(ctx, todoID); err != nil {
		return NewUseCaseError("failed to delete todo", err)
	}

	return nil
}
