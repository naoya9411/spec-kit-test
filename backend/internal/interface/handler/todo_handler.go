package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	todoapp "todo-backend/internal/application/todo"
	"todo-backend/internal/interface/dto"
)

// TodoHandler handles HTTP requests for todo operations
type TodoHandler struct {
	createTodoUseCase *todoapp.CreateTodoUseCase
	getTodosUseCase   *todoapp.GetTodosUseCase
	updateTodoUseCase *todoapp.UpdateTodoUseCase
	deleteTodoUseCase *todoapp.DeleteTodoUseCase
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(
	createTodoUseCase *todoapp.CreateTodoUseCase,
	getTodosUseCase *todoapp.GetTodosUseCase,
	updateTodoUseCase *todoapp.UpdateTodoUseCase,
	deleteTodoUseCase *todoapp.DeleteTodoUseCase,
) *TodoHandler {
	return &TodoHandler{
		createTodoUseCase: createTodoUseCase,
		getTodosUseCase:   getTodosUseCase,
		updateTodoUseCase: updateTodoUseCase,
		deleteTodoUseCase: deleteTodoUseCase,
	}
}

// GetTodos handles GET /api/v1/todos
func (h *TodoHandler) GetTodos(c *gin.Context) {
	var queryParams struct {
		Completed *bool `form:"completed"`
	}

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid query parameters",
			Message: err.Error(),
		})
		return
	}

	request := todoapp.GetTodosRequest{
		Completed: queryParams.Completed,
	}

	response, err := h.getTodosUseCase.Execute(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to retrieve todos",
			Message: err.Error(),
		})
		return
	}

	// Convert to DTO
	todoDTOs := make([]dto.TodoResponse, len(response.Data))
	for i, todo := range response.Data {
		todoDTOs[i] = dto.TodoResponse{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		}
	}

	var statsDTO *dto.CompletionStats
	if response.Stats != nil {
		statsDTO = &dto.CompletionStats{
			Total:     response.Stats.Total,
			Completed: response.Stats.Completed,
			Pending:   response.Stats.Pending,
			Ratio:     response.Stats.Ratio,
		}
	}

	c.JSON(http.StatusOK, dto.GetTodosResponse{
		Data:  todoDTOs,
		Stats: statsDTO,
	})
}

// CreateTodo handles POST /api/v1/todos
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var request dto.CreateTodoRequest
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.ValidationErrorResponse{
			Error:  "Validation failed",
			Fields: map[string]string{"title": "Title is required"},
		})
		return
	}

	useCaseRequest := todoapp.CreateTodoRequest{
		Title:       request.Title,
		Description: request.Description,
	}

	response, err := h.createTodoUseCase.Execute(c.Request.Context(), useCaseRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Failed to create todo",
			Message: err.Error(),
		})
		return
	}

	// Convert to DTO
	todoDTO := dto.TodoResponse{
		ID:          response.ID,
		Title:       response.Title,
		Description: response.Description,
		Completed:   response.Completed,
		CreatedAt:   response.CreatedAt,
		UpdatedAt:   response.UpdatedAt,
	}

	c.JSON(http.StatusCreated, dto.CreateTodoResponse{
		Data: todoDTO,
	})
}

// UpdateTodo handles PUT /api/v1/todos/{id}
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	todoID := c.Param("id")
	
	var request dto.UpdateTodoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	useCaseRequest := todoapp.UpdateTodoRequest{
		ID:          todoID,
		Title:       request.Title,
		Description: request.Description,
		Completed:   request.Completed,
	}

	response, err := h.updateTodoUseCase.Execute(c.Request.Context(), useCaseRequest)
	if err != nil {
		// Check if it's a not found error
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Todo not found",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Failed to update todo",
			Message: err.Error(),
		})
		return
	}

	// Convert to DTO
	todoDTO := dto.TodoResponse{
		ID:          response.ID,
		Title:       response.Title,
		Description: response.Description,
		Completed:   response.Completed,
		CreatedAt:   response.CreatedAt,
		UpdatedAt:   response.UpdatedAt,
	}

	c.JSON(http.StatusOK, dto.UpdateTodoResponse{
		Data: todoDTO,
	})
}

// DeleteTodo handles DELETE /api/v1/todos/{id}
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	todoID := c.Param("id")

	request := todoapp.DeleteTodoRequest{
		ID: todoID,
	}

	err := h.deleteTodoUseCase.Execute(c.Request.Context(), request)
	if err != nil {
		// Check if it's a not found error
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Todo not found",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to delete todo",
			Message: err.Error(),
		})
		return
	}

	// Return 204 No Content for successful deletion
	c.Status(http.StatusNoContent)
}

// contains is a utility function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		   len(s) > len(substr) && s[len(s)-len(substr):] == substr ||
		   (len(s) > len(substr) && func() bool {
			   for i := 1; i < len(s)-len(substr)+1; i++ {
				   if s[i:i+len(substr)] == substr {
					   return true
				   }
			   }
			   return false
		   }())
}
