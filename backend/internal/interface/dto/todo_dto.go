package dto

// CreateTodoRequest represents the request body for creating a todo
type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required" example:"Buy groceries"`
	Description string `json:"description,omitempty" example:"Buy milk, eggs, and bread"`
}

// UpdateTodoRequest represents the request body for updating a todo
type UpdateTodoRequest struct {
	Title       *string `json:"title,omitempty" example:"Updated todo title"`
	Description *string `json:"description,omitempty" example:"Updated description"`
	Completed   *bool   `json:"completed,omitempty" example:"true"`
}

// TodoResponse represents a todo in the response
type TodoResponse struct {
	ID          string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Title       string `json:"title" example:"Buy groceries"`
	Description string `json:"description,omitempty" example:"Buy milk, eggs, and bread"`
	Completed   bool   `json:"completed" example:"false"`
	CreatedAt   string `json:"created_at" example:"2023-01-01T10:00:00Z"`
	UpdatedAt   string `json:"updated_at" example:"2023-01-01T10:00:00Z"`
}

// GetTodosResponse represents the response for getting multiple todos
type GetTodosResponse struct {
	Data  []TodoResponse  `json:"data"`
	Stats *CompletionStats `json:"stats,omitempty"`
}

// CreateTodoResponse represents the response after creating a todo
type CreateTodoResponse struct {
	Data TodoResponse `json:"data"`
}

// UpdateTodoResponse represents the response after updating a todo
type UpdateTodoResponse struct {
	Data TodoResponse `json:"data"`
}

// CompletionStats represents statistics about todo completion
type CompletionStats struct {
	Total     int     `json:"total" example:"10"`
	Completed int     `json:"completed" example:"7"`
	Pending   int     `json:"pending" example:"3"`
	Ratio     float64 `json:"completion_ratio" example:"0.7"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Validation failed"`
	Message string `json:"message,omitempty" example:"Title is required"`
}

// ValidationErrorResponse represents a validation error response with field details
type ValidationErrorResponse struct {
	Error  string            `json:"error" example:"Validation failed"`
	Fields map[string]string `json:"fields" example:"{\"title\": \"Title is required\"}"`
}
