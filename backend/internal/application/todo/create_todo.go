package todo

import (
	"context"

	"todo-backend/internal/domain/todo"
)

// CreateTodoRequest represents the input for creating a todo
type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description,omitempty"`
}

// CreateTodoResponse represents the output after creating a todo
type CreateTodoResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateTodoUseCase handles the creation of new todos
type CreateTodoUseCase struct {
	repository     todo.Repository
	domainService  todo.Service
}

// NewCreateTodoUseCase creates a new CreateTodoUseCase
func NewCreateTodoUseCase(repository todo.Repository, domainService todo.Service) *CreateTodoUseCase {
	return &CreateTodoUseCase{
		repository:    repository,
		domainService: domainService,
	}
}

// Execute creates a new todo
func (uc *CreateTodoUseCase) Execute(ctx context.Context, request CreateTodoRequest) (*CreateTodoResponse, error) {
	// Validate the request using domain service
	if err := uc.domainService.ValidateCreate(ctx, request.Title, request.Description); err != nil {
		return nil, NewUseCaseError("validation failed", err)
	}

	// Create title value object
	title, err := todo.NewTitle(request.Title)
	if err != nil {
		return nil, NewUseCaseError("invalid title", err)
	}

	// Create new todo entity
	newTodo := todo.NewTodo(title, request.Description)

	// Save to repository
	if err := uc.repository.Create(ctx, newTodo); err != nil {
		return nil, NewUseCaseError("failed to create todo", err)
	}

	// Convert to response
	response := &CreateTodoResponse{
		ID:          newTodo.ID.String(),
		Title:       newTodo.Title.String(),
		Description: newTodo.Description,
		Completed:   newTodo.Completed,
		CreatedAt:   newTodo.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   newTodo.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return response, nil
}
