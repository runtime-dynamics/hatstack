package testutil

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetupTestRouter(t *testing.T) {
	router := SetupTestRouter()

	if router == nil {
		t.Fatal("SetupTestRouter() returned nil")
	}

	// Verify Gin is in test mode
	if gin.Mode() != gin.TestMode {
		t.Error("Gin should be in test mode")
	}
}

func TestCreateTestContext(t *testing.T) {
	w, c := CreateTestContext()

	if w == nil {
		t.Fatal("CreateTestContext() returned nil response recorder")
	}

	if c == nil {
		t.Fatal("CreateTestContext() returned nil context")
	}

	// Verify we can write to the context
	c.String(200, "test")
	if w.Body.String() != "test" {
		t.Error("Failed to write to test context")
	}
}

func TestAssertNoError(t *testing.T) {
	// This test verifies AssertNoError doesn't fail with nil error
	AssertNoError(t, nil, "should not fail")
}

func TestAssertError(t *testing.T) {
	// This test verifies AssertError doesn't fail with non-nil error
	err := errors.New("test error")
	AssertError(t, err, "should not fail")
}
