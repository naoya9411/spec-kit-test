package todo

import (
	"context"
	"errors"

	"todo-backend/internal/domain/todo"
)

// UpdateTodoRequest represents the input for updating a todo
type UpdateTodoRequest struct {
	ID          string  `uri:"id" binding:"required"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
}

// UpdateTodoResponse represents the output after updating a todo
type UpdateTodoResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// UpdateTodoUseCase handles updating existing todos
type UpdateTodoUseCase struct {
	repository     todo.Repository
	domainService  todo.Service
}

// NewUpdateTodoUseCase creates a new UpdateTodoUseCase
func NewUpdateTodoUseCase(repository todo.Repository, domainService todo.Service) *UpdateTodoUseCase {
	return &UpdateTodoUseCase{
		repository:    repository,
		domainService: domainService,
	}
}

// Execute updates an existing todo
func (uc *UpdateTodoUseCase) Execute(ctx context.Context, request UpdateTodoRequest) (*UpdateTodoResponse, error) {
	// Parse todo ID
	todoID, err := todo.NewTodoIDFromString(request.ID)
	if err != nil {
		return nil, NewUseCaseError("invalid todo ID", err)
	}

	// Get existing todo
	existingTodo, err := uc.repository.GetByID(ctx, todoID)
	if err != nil {
		return nil, NewUseCaseError("failed to get todo", err)
	}

	if existingTodo == nil {
		return nil, NewUseCaseError("todo not found", errors.New("todo does not exist"))
	}

	// Validate update request
	if err := uc.domainService.ValidateUpdate(ctx, existingTodo, request.Title, request.Description, request.Completed); err != nil {
		return nil, NewUseCaseError("validation failed", err)
	}

	// Apply updates
	if request.Title != nil {
		newTitle, err := todo.NewTitle(*request.Title)
		if err != nil {
			return nil, NewUseCaseError("invalid title", err)
		}
		existingTodo.UpdateTitle(newTitle)
	}

	if request.Description != nil {
		existingTodo.UpdateDescription(*request.Description)
	}

	if request.Completed != nil {
		if *request.Completed {
			existingTodo.MarkAsCompleted()
		} else {
			existingTodo.MarkAsIncomplete()
		}
	}

	// Save updated todo
	if err := uc.repository.Update(ctx, existingTodo); err != nil {
		return nil, NewUseCaseError("failed to update todo", err)
	}

	// Convert to response
	response := &UpdateTodoResponse{
		ID:          existingTodo.ID.String(),
		Title:       existingTodo.Title.String(),
		Description: existingTodo.Description,
		Completed:   existingTodo.Completed,
		CreatedAt:   existingTodo.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   existingTodo.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return response, nil
}
