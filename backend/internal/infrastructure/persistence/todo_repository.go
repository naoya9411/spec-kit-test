package persistence

import (
	"context"
	"errors"

	"todo-backend/internal/domain/todo"

	"gorm.io/gorm"
)

// todoRepository implements the todo.Repository interface using GORM
type todoRepository struct {
	db *gorm.DB
}

// NewTodoRepository creates a new todo repository
func NewTodoRepository(db *gorm.DB) todo.Repository {
	return &todoRepository{db: db}
}

// Create saves a new todo to the repository
func (r *todoRepository) Create(ctx context.Context, t *todo.Todo) error {
	if err := r.db.WithContext(ctx).Create(t).Error; err != nil {
		return todo.NewRepositoryError("create", err)
	}
	return nil
}

// GetByID retrieves a todo by its ID
func (r *todoRepository) GetByID(ctx context.Context, id todo.TodoID) (*todo.Todo, error) {
	var t todo.Todo
	err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&t).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil for not found instead of error
		}
		return nil, todo.NewRepositoryError("get by ID", err)
	}
	
	return &t, nil
}

// GetAll retrieves all todos
func (r *todoRepository) GetAll(ctx context.Context) ([]*todo.Todo, error) {
	var todos []*todo.Todo
	if err := r.db.WithContext(ctx).Find(&todos).Error; err != nil {
		return nil, todo.NewRepositoryError("get all", err)
	}
	return todos, nil
}

// Update updates an existing todo
func (r *todoRepository) Update(ctx context.Context, t *todo.Todo) error {
	result := r.db.WithContext(ctx).Save(t)
	if result.Error != nil {
		return todo.NewRepositoryError("update", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return todo.NewRepositoryError("update", errors.New("no rows affected"))
	}
	
	return nil
}

// Delete removes a todo by its ID
func (r *todoRepository) Delete(ctx context.Context, id todo.TodoID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id.String()).Delete(&todo.Todo{})
	if result.Error != nil {
		return todo.NewRepositoryError("delete", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return todo.NewRepositoryError("delete", errors.New("no rows affected"))
	}
	
	return nil
}

// FindByCompleted retrieves todos filtered by completion status
func (r *todoRepository) FindByCompleted(ctx context.Context, completed bool) ([]*todo.Todo, error) {
	var todos []*todo.Todo
	if err := r.db.WithContext(ctx).Where("completed = ?", completed).Find(&todos).Error; err != nil {
		return nil, todo.NewRepositoryError("find by completed", err)
	}
	return todos, nil
}

// Exists checks if a todo exists by its ID
func (r *todoRepository) Exists(ctx context.Context, id todo.TodoID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&todo.Todo{}).Where("id = ?", id.String()).Count(&count).Error; err != nil {
		return false, todo.NewRepositoryError("exists", err)
	}
	return count > 0, nil
}
