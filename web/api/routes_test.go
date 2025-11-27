package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
}

func TestRenderError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		statusCode     int
		message        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "not found error",
			err:            datastore.ErrNoSuchEntity,
			statusCode:     http.StatusInternalServerError,
			message:        "entity not found",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"entity not found"}`,
		},
		{
			name:           "generic error",
			err:            errors.New("generic error"),
			statusCode:     http.StatusInternalServerError,
			message:        "internal server error",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
		{
			name:           "bad request error",
			err:            errors.New("bad request"),
			statusCode:     http.StatusBadRequest,
			message:        "invalid input",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid input"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			renderError(c, tt.err, tt.statusCode, tt.message)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestRenderSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	renderSuccess(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message":"success"}`, w.Body.String())
}

func TestRenderFinal(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		message        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success case",
			err:            nil,
			message:        "operation completed",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"success"}`,
		},
		{
			name:           "error case",
			err:            errors.New("something went wrong"),
			message:        "operation failed",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"operation failed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			renderFinal(c, tt.err, tt.message)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestRenderFinalContent(t *testing.T) {
	tests := []struct {
		name           string
		content        interface{}
		key            string
		err            error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success with key",
			content:        map[string]string{"name": "test"},
			key:            "data",
			err:            nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":{"name":"test"}}`,
		},
		{
			name:           "success without key",
			content:        map[string]string{"name": "test"},
			key:            "",
			err:            nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"name":"test"}`,
		},
		{
			name:           "error case",
			content:        nil,
			key:            "data",
			err:            errors.New("failed to fetch"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"data"}`,
		},
		{
			name:           "success with array",
			content:        []string{"item1", "item2"},
			key:            "items",
			err:            nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"items":["item1","item2"]}`,
		},
		{
			name:           "success with string",
			content:        "simple string",
			key:            "message",
			err:            nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"simple string"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			renderFinalContent(c, tt.content, tt.key, tt.err)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestRegisterRoutes(t *testing.T) {
	router := gin.New()
	
	// Should not panic
	RegisterRoutes(router)

	// Verify router was created successfully
	if router == nil {
		t.Error("Router should not be nil after RegisterRoutes")
	}
}

// Benchmark tests
func BenchmarkRenderSuccess(b *testing.B) {
	gin.SetMode(gin.TestMode)
	
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		renderSuccess(c)
	}
}

func BenchmarkRenderError(b *testing.B) {
	gin.SetMode(gin.TestMode)
	err := errors.New("test error")
	
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		renderError(c, err, http.StatusInternalServerError, "test message")
	}
}

func BenchmarkRenderFinalContent(b *testing.B) {
	gin.SetMode(gin.TestMode)
	content := map[string]string{"test": "data"}
	
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		renderFinalContent(c, content, "data", nil)
	}
}
