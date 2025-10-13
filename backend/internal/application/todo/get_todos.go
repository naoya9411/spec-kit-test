package todo

import (
	"context"

	"todo-backend/internal/domain/todo"
)

// GetTodosRequest represents the input for getting todos
type GetTodosRequest struct {
	Completed *bool `form:"completed,omitempty"`
}

// GetTodosResponse represents the output after getting todos
type GetTodosResponse struct {
	Data  []TodoDTO           `json:"data"`
	Stats *todo.CompletionStats `json:"stats,omitempty"`
}

// TodoDTO represents a todo in the response
type TodoDTO struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// GetTodosUseCase handles retrieving todos
type GetTodosUseCase struct {
	repository     todo.Repository
	domainService  todo.Service
}

// NewGetTodosUseCase creates a new GetTodosUseCase
func NewGetTodosUseCase(repository todo.Repository, domainService todo.Service) *GetTodosUseCase {
	return &GetTodosUseCase{
		repository:    repository,
		domainService: domainService,
	}
}

// Execute retrieves todos based on the request parameters
func (uc *GetTodosUseCase) Execute(ctx context.Context, request GetTodosRequest) (*GetTodosResponse, error) {
	var todos []*todo.Todo
	var err error

	// Filter by completion status if specified
	if request.Completed != nil {
		todos, err = uc.repository.FindByCompleted(ctx, *request.Completed)
	} else {
		todos, err = uc.repository.GetAll(ctx)
	}

	if err != nil {
		return nil, NewUseCaseError("failed to retrieve todos", err)
	}

	// Convert to DTOs
	todoDTOs := make([]TodoDTO, len(todos))
	for i, t := range todos {
		todoDTOs[i] = TodoDTO{
			ID:          t.ID.String(),
			Title:       t.Title.String(),
			Description: t.Description,
			Completed:   t.Completed,
			CreatedAt:   t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   t.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	// Get completion statistics
	stats := uc.domainService.GetCompletionStats(ctx, todos)

	response := &GetTodosResponse{
		Data:  todoDTOs,
		Stats: stats,
	}

	return response, nil
}
