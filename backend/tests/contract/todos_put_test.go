package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestPutTodos tests the PUT /api/v1/todos/{id} endpoint contract
func TestPutTodos(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name           string
		todoID         string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedFields []string
	}{
		{
			name:   "should update todo and return 200",
			todoID: "1",
			requestBody: map[string]interface{}{
				"title":       "Updated Todo",
				"description": "Updated Description",
				"completed":   true,
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"data"},
		},
		{
			name:   "should return 404 for non-existent todo",
			todoID: "999",
			requestBody: map[string]interface{}{
				"title": "Updated Todo",
			},
			expectedStatus: http.StatusNotFound,
			expectedFields: []string{"error"},
		},
		{
			name:   "should return 400 for invalid request body",
			todoID: "1",
			requestBody: map[string]interface{}{
				"title": "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router (this will fail until we implement the handler)
			router := setupTestRouter()

			// Create request body
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("PUT", "/api/v1/todos/"+tt.todoID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert response structure
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert required fields exist
			for _, field := range tt.expectedFields {
				assert.Contains(t, response, field)
			}

			// If successful update, assert the updated todo structure
			if tt.expectedStatus == http.StatusOK {
				if data, exists := response["data"]; exists {
					todoData := data.(map[string]interface{})
					assert.Contains(t, todoData, "id")
					assert.Contains(t, todoData, "title")
					assert.Contains(t, todoData, "completed")
					assert.Contains(t, todoData, "updated_at")
				}
			}
		})
	}
}
