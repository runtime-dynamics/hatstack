package services

import (
	"context"
	"testing"
	"time"
)

func TestNewBaseService(t *testing.T) {
	ctx := context.Background()
	service := NewBaseService(ctx)

	if service == nil {
		t.Fatal("NewBaseService() returned nil")
	}

	if service.ctx != ctx {
		t.Error("NewBaseService() did not set context correctly")
	}
}

func TestBaseService_Context(t *testing.T) {
	ctx := context.Background()
	service := NewBaseService(ctx)

	returnedCtx := service.Context()

	if returnedCtx != ctx {
		t.Error("Context() returned different context than was set")
	}
}

func TestBaseService_WithCancelContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := NewBaseService(ctx)

	// Verify context is not cancelled initially
	select {
	case <-service.Context().Done():
		t.Error("Context should not be cancelled initially")
	default:
		// Expected
	}

	// Cancel the context
	cancel()

	// Verify context is now cancelled
	select {
	case <-service.Context().Done():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should be cancelled after cancel() is called")
	}
}

func TestBaseService_WithTimeoutContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	service := NewBaseService(ctx)

	// Wait for timeout
	select {
	case <-service.Context().Done():
		// Expected - context should timeout
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should have timed out")
	}
}

func TestBaseService_WithValueContext(t *testing.T) {
	type key string
	const userKey key = "user"

	ctx := context.WithValue(context.Background(), userKey, "testuser")
	service := NewBaseService(ctx)

	// Verify we can retrieve the value from the service's context
	value := service.Context().Value(userKey)
	if value == nil {
		t.Fatal("Context value should not be nil")
	}

	if value.(string) != "testuser" {
		t.Errorf("Context value = %v, want %v", value, "testuser")
	}
}

func TestBaseService_ImplementsServiceInterface(t *testing.T) {
	ctx := context.Background()
	service := NewBaseService(ctx)

	// Verify BaseService implements Service interface
	var _ Service = service
}

func TestBaseService_MultipleInstances(t *testing.T) {
	type contextKey string
	const testKey contextKey = "key"

	ctx1 := context.WithValue(context.Background(), testKey, "value1")
	ctx2 := context.WithValue(context.Background(), testKey, "value2")

	service1 := NewBaseService(ctx1)
	service2 := NewBaseService(ctx2)

	// Verify each service has its own context
	if service1.Context() == service2.Context() {
		t.Error("Different service instances should have different contexts")
	}

	// Verify values are different
	val1 := service1.Context().Value(testKey)
	val2 := service2.Context().Value(testKey)

	if val1 == val2 {
		t.Error("Different service instances should have different context values")
	}

	if val1.(string) != "value1" {
		t.Errorf("service1 context value = %v, want %v", val1, "value1")
	}

	if val2.(string) != "value2" {
		t.Errorf("service2 context value = %v, want %v", val2, "value2")
	}
}
