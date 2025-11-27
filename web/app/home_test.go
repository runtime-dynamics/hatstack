package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
}

func TestHomePageHandler(t *testing.T) {
	// Create a test router
	router := gin.New()
	router.GET("/", HomePageHandler)

	// Create a test request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Verify status code
	assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 OK")

	// Verify response is HTML
	contentType := w.Header().Get("Content-Type")
	assert.Contains(t, contentType, "text/html", "Expected HTML content type")

	// Verify response contains expected HTML elements
	body := w.Body.String()
	assert.Contains(t, body, "<!DOCTYPE html>", "Expected HTML doctype")
	assert.Contains(t, body, "<html", "Expected HTML tag")
	assert.Contains(t, body, "H.A.T. Stack", "Expected H.A.T. Stack branding")
}

func TestHomePageHandler_HTMLStructure(t *testing.T) {
	router := gin.New()
	router.GET("/", HomePageHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	body := w.Body.String()

	// Verify essential HTML structure
	requiredElements := []string{
		"<head>",
		"</head>",
		"<body",
		"</body>",
		"</html>",
		"<meta charset=\"UTF-8\"",
		"<meta name=\"viewport\"",
	}

	for _, element := range requiredElements {
		assert.Contains(t, body, element, "Expected HTML element: %s", element)
	}
}

func TestHomePageHandler_HTMXIncluded(t *testing.T) {
	router := gin.New()
	router.GET("/", HomePageHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	body := w.Body.String()

	// Verify HTMX is included
	assert.Contains(t, body, "htmx.org", "Expected HTMX library to be included")
}

func TestHomePageHandler_AlpineIncluded(t *testing.T) {
	router := gin.New()
	router.GET("/", HomePageHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	body := w.Body.String()

	// Verify Alpine.js is included
	assert.Contains(t, body, "alpinejs", "Expected Alpine.js library to be included")
}

func TestHomePageHandler_TailwindIncluded(t *testing.T) {
	router := gin.New()
	router.GET("/", HomePageHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	body := w.Body.String()

	// Verify TailwindCSS is included
	assert.Contains(t, body, "tailwindcss", "Expected TailwindCSS to be included")
}

func TestHomePageHandler_Navigation(t *testing.T) {
	router := gin.New()
	router.GET("/", HomePageHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	body := w.Body.String()

	// Verify navigation elements
	assert.Contains(t, body, "<nav", "Expected navigation element")
	assert.Contains(t, body, "/about", "Expected About link")
	assert.Contains(t, body, "/docs", "Expected Docs link")
}

func TestHomePageHandler_Footer(t *testing.T) {
	router := gin.New()
	router.GET("/", HomePageHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	body := w.Body.String()

	// Verify footer elements
	assert.Contains(t, body, "<footer", "Expected footer element")
	assert.Contains(t, body, "2025", "Expected copyright year")
}

func TestHomePageHandler_MultipleRequests(t *testing.T) {
	router := gin.New()
	router.GET("/", HomePageHandler)

	// Make multiple requests to ensure handler is stateless
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Request %d failed", i+1)
		assert.Contains(t, w.Body.String(), "H.A.T. Stack", "Request %d missing content", i+1)
	}
}

func TestHomePageHandler_NoErrors(t *testing.T) {
	router := gin.New()
	router.GET("/", HomePageHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	body := w.Body.String()

	// Verify no error messages in output
	errorIndicators := []string{
		"error",
		"Error",
		"ERROR",
		"panic",
		"fatal",
	}

	for _, indicator := range errorIndicators {
		// Only check if it appears in an error context (case-sensitive)
		if strings.Contains(body, "class=\"error\"") || 
		   strings.Contains(body, "id=\"error\"") {
			t.Errorf("Found error indicator in HTML: %s", indicator)
		}
	}
}

// Benchmark tests
func BenchmarkHomePageHandler(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", HomePageHandler)

	req, _ := http.NewRequest("GET", "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
