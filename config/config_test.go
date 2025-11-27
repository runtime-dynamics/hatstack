package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Save original env vars
	originalDatastore := os.Getenv("DATASTORE_NAME")
	originalFrontend := os.Getenv("FRONTEND_ENDPOINT")
	originalFirebaseKey := os.Getenv("FIREBASE_API_KEY")
	originalFirebaseDomain := os.Getenv("FIREBASE_AUTH_DOMAIN")

	// Restore env vars after test
	defer func() {
		os.Setenv("DATASTORE_NAME", originalDatastore)
		os.Setenv("FRONTEND_ENDPOINT", originalFrontend)
		os.Setenv("FIREBASE_API_KEY", originalFirebaseKey)
		os.Setenv("FIREBASE_AUTH_DOMAIN", originalFirebaseDomain)
	}()

	tests := []struct {
		name                string
		envVars             map[string]string
		expectedDatastore   string
		expectedFrontend    string
		expectedFirebaseKey string
		expectedFirebaseDom string
	}{
		{
			name: "all env vars set",
			envVars: map[string]string{
				"DATASTORE_NAME":       "test-datastore",
				"FRONTEND_ENDPOINT":    "http://test.example.com",
				"FIREBASE_API_KEY":     "test-api-key",
				"FIREBASE_AUTH_DOMAIN": "test.firebaseapp.com",
			},
			expectedDatastore:   "test-datastore",
			expectedFrontend:    "http://test.example.com",
			expectedFirebaseKey: "test-api-key",
			expectedFirebaseDom: "test.firebaseapp.com",
		},
		{
			name:                "no env vars set - uses defaults",
			envVars:             map[string]string{},
			expectedDatastore:   "default",
			expectedFrontend:    "http://local.nitecon.net:8080",
			expectedFirebaseKey: "",
			expectedFirebaseDom: "",
		},
		{
			name: "partial env vars set",
			envVars: map[string]string{
				"DATASTORE_NAME":    "partial-datastore",
				"FRONTEND_ENDPOINT": "http://partial.example.com",
			},
			expectedDatastore:   "partial-datastore",
			expectedFrontend:    "http://partial.example.com",
			expectedFirebaseKey: "",
			expectedFirebaseDom: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env vars
			os.Unsetenv("DATASTORE_NAME")
			os.Unsetenv("FRONTEND_ENDPOINT")
			os.Unsetenv("FIREBASE_API_KEY")
			os.Unsetenv("FIREBASE_AUTH_DOMAIN")

			// Set test env vars
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Load config
			err := LoadConfig()
			if err != nil {
				t.Fatalf("LoadConfig() error = %v", err)
			}

			// Get config
			cfg := Get()

			// Verify values
			if cfg.DataStoreName != tt.expectedDatastore {
				t.Errorf("DataStoreName = %v, want %v", cfg.DataStoreName, tt.expectedDatastore)
			}
			if cfg.FrontendEndpoint != tt.expectedFrontend {
				t.Errorf("FrontendEndpoint = %v, want %v", cfg.FrontendEndpoint, tt.expectedFrontend)
			}
			if cfg.FirebaseAPIKey != tt.expectedFirebaseKey {
				t.Errorf("FirebaseAPIKey = %v, want %v", cfg.FirebaseAPIKey, tt.expectedFirebaseKey)
			}
			if cfg.FirebaseAuthDomain != tt.expectedFirebaseDom {
				t.Errorf("FirebaseAuthDomain = %v, want %v", cfg.FirebaseAuthDomain, tt.expectedFirebaseDom)
			}
		})
	}
}

func TestGet(t *testing.T) {
	// Load a test config
	os.Setenv("DATASTORE_NAME", "test-get")
	defer os.Unsetenv("DATASTORE_NAME")

	err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Get config multiple times to test thread safety
	cfg1 := Get()
	cfg2 := Get()

	if cfg1 == nil {
		t.Error("Get() returned nil")
	}
	if cfg2 == nil {
		t.Error("Get() returned nil on second call")
	}
	if cfg1 != cfg2 {
		t.Error("Get() returned different instances")
	}
}

func TestDecodeBase64Cert(t *testing.T) {
	tests := []struct {
		name     string
		envVar   string
		envValue string
		expected string
	}{
		{
			name:     "empty env var",
			envVar:   "TEST_CERT_EMPTY",
			envValue: "",
			expected: "",
		},
		{
			name:     "plain text (not base64)",
			envVar:   "TEST_CERT_PLAIN",
			envValue: "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
			expected: "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
		},
		{
			name:     "base64 encoded",
			envVar:   "TEST_CERT_BASE64",
			envValue: "dGVzdCBjZXJ0aWZpY2F0ZQ==", // "test certificate" in base64
			expected: "test certificate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.envVar, tt.envValue)
			defer os.Unsetenv(tt.envVar)

			result := decodeBase64Cert(tt.envVar)
			if result != tt.expected {
				t.Errorf("decodeBase64Cert() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConfigThreadSafety(t *testing.T) {
	// Test concurrent access to config
	os.Setenv("DATASTORE_NAME", "thread-test")
	defer os.Unsetenv("DATASTORE_NAME")

	err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Run multiple goroutines accessing config
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			cfg := Get()
			if cfg == nil {
				t.Error("Get() returned nil in goroutine")
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}
