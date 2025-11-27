package data

import (
	"errors"
	"testing"

	"cloud.google.com/go/datastore"
)

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "datastore.ErrNoSuchEntity",
			err:      datastore.ErrNoSuchEntity,
			expected: true,
		},
		{
			name:     "wrapped datastore.ErrNoSuchEntity",
			err:      errors.New("wrapped: " + datastore.ErrNoSuchEntity.Error()),
			expected: false, // errors.Is won't match wrapped string
		},
		{
			name:     "different error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotFound(tt.err)
			if result != tt.expected {
				t.Errorf("IsNotFound(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

func TestIsNotFound_WithErrorsIs(t *testing.T) {
	// Test that IsNotFound works with errors.Is
	wrappedErr := errors.Join(errors.New("context"), datastore.ErrNoSuchEntity)
	
	if !IsNotFound(wrappedErr) {
		t.Error("IsNotFound should return true for wrapped ErrNoSuchEntity using errors.Join")
	}
}

// Note: Testing Cli() requires actual Google Cloud credentials and is better suited
// for integration tests. Here we document what should be tested in integration:
//
// Integration tests for Cli() should verify:
// - Client is created successfully with valid credentials
// - Client is reused on subsequent calls (singleton pattern)
// - Client fails gracefully with invalid credentials
// - Thread-safe access to the client
//
// Example integration test structure (not run in unit tests):
/*
func TestCli_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Set up test environment with valid credentials
	os.Setenv("GOOGLE_PROJECT_ID", "test-project")
	os.Setenv("DATASTORE_NAME", "test-database")
	defer func() {
		os.Unsetenv("GOOGLE_PROJECT_ID")
		os.Unsetenv("DATASTORE_NAME")
	}()

	// Test client creation
	client1 := Cli()
	if client1 == nil {
		t.Fatal("Cli() returned nil")
	}

	// Test singleton pattern
	client2 := Cli()
	if client1 != client2 {
		t.Error("Cli() should return the same client instance")
	}
}
*/
