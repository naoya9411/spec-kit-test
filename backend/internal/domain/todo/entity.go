package todo

import (
	"time"
)

// Todo represents a todo item in the domain
type Todo struct {
	ID          TodoID    `json:"id" gorm:"primaryKey"`
	Title       Title     `json:"title" gorm:"type:varchar(255);not null"`
	Description string    `json:"description,omitempty" gorm:"type:text"`
	Completed   bool      `json:"completed" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// NewTodo creates a new Todo with validation
func NewTodo(title Title, description string) *Todo {
	return &Todo{
		ID:          NewTodoID(),
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// MarkAsCompleted marks the todo as completed
func (t *Todo) MarkAsCompleted() {
	t.Completed = true
	t.UpdatedAt = time.Now()
}

// MarkAsIncomplete marks the todo as incomplete
func (t *Todo) MarkAsIncomplete() {
	t.Completed = false
	t.UpdatedAt = time.Now()
}

// UpdateTitle updates the todo title
func (t *Todo) UpdateTitle(title Title) {
	t.Title = title
	t.UpdatedAt = time.Now()
}

// UpdateDescription updates the todo description
func (t *Todo) UpdateDescription(description string) {
	t.Description = description
	t.UpdatedAt = time.Now()
}

// IsCompleted returns whether the todo is completed
func (t *Todo) IsCompleted() bool {
	return t.Completed
}

// Validate validates the todo entity
func (t *Todo) Validate() error {
	if err := t.Title.Validate(); err != nil {
		return err
	}
	return nil
}
