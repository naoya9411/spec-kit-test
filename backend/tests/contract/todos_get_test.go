package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGetTodos tests the GET /api/v1/todos endpoint contract
func TestGetTodos(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name           string
		expectedStatus int
		expectedFields []string
	}{
		{
			name:           "should return 200 and list of todos",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"data"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router (this will fail until we implement the handler)
			router := setupTestRouter()

			// Create request
			req, _ := http.NewRequest("GET", "/api/v1/todos", nil)
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

			// Assert data is an array
			if data, exists := response["data"]; exists {
				assert.IsType(t, []interface{}{}, data)
			}
		})
	}
}

// setupTestRouter creates a test router - this will be implemented later
func setupTestRouter() *gin.Engine {
	// This is a placeholder that will fail until we implement the actual router
	router := gin.New()
	
	// TODO: Add actual routes when handlers are implemented
	// For now, this will cause the test to fail as expected in TDD
	
	return router
}
