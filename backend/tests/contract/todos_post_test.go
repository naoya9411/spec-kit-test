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

// TestPostTodos tests the POST /api/v1/todos endpoint contract
func TestPostTodos(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedFields []string
	}{
		{
			name: "should create todo and return 201",
			requestBody: map[string]interface{}{
				"title":       "Test Todo",
				"description": "Test Description",
			},
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"data"},
		},
		{
			name: "should return 400 for missing title",
			requestBody: map[string]interface{}{
				"description": "Test Description without title",
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
			req, _ := http.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(body))
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

			// If successful creation, assert the created todo structure
			if tt.expectedStatus == http.StatusCreated {
				if data, exists := response["data"]; exists {
					todoData := data.(map[string]interface{})
					assert.Contains(t, todoData, "id")
					assert.Contains(t, todoData, "title")
					assert.Contains(t, todoData, "completed")
					assert.Contains(t, todoData, "created_at")
				}
			}
		})
	}
}
