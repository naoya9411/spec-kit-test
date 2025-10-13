package todo

import (
	"context"
	"errors"
)

// Service defines the domain service contract for todo operations
type Service interface {
	// ValidateCreate validates a todo for creation
	ValidateCreate(ctx context.Context, title, description string) error

	// ValidateUpdate validates a todo for update
	ValidateUpdate(ctx context.Context, todo *Todo, title, description *string, completed *bool) error

	// CanDelete checks if a todo can be deleted
	CanDelete(ctx context.Context, id TodoID) error

	// GetCompletionStats returns statistics about todo completion
	GetCompletionStats(ctx context.Context, todos []*Todo) *CompletionStats
}

// CompletionStats represents statistics about todo completion
type CompletionStats struct {
	Total     int     `json:"total"`
	Completed int     `json:"completed"`
	Pending   int     `json:"pending"`
	Ratio     float64 `json:"completion_ratio"`
}

// DomainService implements the domain service
type DomainService struct {
	repository Repository
}

// NewDomainService creates a new domain service
func NewDomainService(repository Repository) Service {
	return &DomainService{
		repository: repository,
	}
}

// ValidateCreate validates a todo for creation
func (s *DomainService) ValidateCreate(ctx context.Context, title, description string) error {
	// Validate title
	if _, err := NewTitle(title); err != nil {
		return err
	}

	// Validate description length
	if len(description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}

	return nil
}

// ValidateUpdate validates a todo for update
func (s *DomainService) ValidateUpdate(ctx context.Context, todo *Todo, title, description *string, completed *bool) error {
	// Validate title if provided
	if title != nil {
		if _, err := NewTitle(*title); err != nil {
			return err
		}
	}

	// Validate description if provided
	if description != nil && len(*description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}

	return nil
}

// CanDelete checks if a todo can be deleted
func (s *DomainService) CanDelete(ctx context.Context, id TodoID) error {
	// Check if todo exists
	exists, err := s.repository.Exists(ctx, id)
	if err != nil {
		return NewRepositoryError("check existence", err)
	}

	if !exists {
		return errors.New("todo not found")
	}

	// In this simple domain, all todos can be deleted
	// In more complex scenarios, you might have business rules here
	return nil
}

// GetCompletionStats returns statistics about todo completion
func (s *DomainService) GetCompletionStats(ctx context.Context, todos []*Todo) *CompletionStats {
	if len(todos) == 0 {
		return &CompletionStats{
			Total:     0,
			Completed: 0,
			Pending:   0,
			Ratio:     0,
		}
	}

	stats := &CompletionStats{
		Total: len(todos),
	}

	for _, todo := range todos {
		if todo.IsCompleted() {
			stats.Completed++
		} else {
			stats.Pending++
		}
	}

	if stats.Total > 0 {
		stats.Ratio = float64(stats.Completed) / float64(stats.Total)
	}

	return stats
}
