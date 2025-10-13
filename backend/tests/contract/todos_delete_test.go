package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestDeleteTodos tests the DELETE /api/v1/todos/{id} endpoint contract
func TestDeleteTodos(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name           string
		todoID         string
		expectedStatus int
		expectedFields []string
	}{
		{
			name:           "should delete todo and return 204",
			todoID:         "1",
			expectedStatus: http.StatusNoContent,
			expectedFields: []string{}, // No content expected for 204
		},
		{
			name:           "should return 404 for non-existent todo",
			todoID:         "999",
			expectedStatus: http.StatusNotFound,
			expectedFields: []string{"error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router (this will fail until we implement the handler)
			router := setupTestRouter()

			// Create request
			req, _ := http.NewRequest("DELETE", "/api/v1/todos/"+tt.todoID, nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// For 204 No Content, there should be no response body
			if tt.expectedStatus == http.StatusNoContent {
				assert.Empty(t, w.Body.String())
			} else {
				// Assert response structure for error cases
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Assert required fields exist
				for _, field := range tt.expectedFields {
					assert.Contains(t, response, field)
				}
			}
		})
	}
}
