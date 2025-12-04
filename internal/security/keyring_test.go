package security

import (
	"testing"

	"github.com/zalando/go-keyring"
)

func TestKeyring(t *testing.T) {
	// Use the mock keyring for testing to avoid OS dependencies in CI/Headless envs
	keyring.MockInit()

	provider := "TestProvider"
	key := "sk-test-12345"

	// Test SetAPIKey
	err := SetAPIKey(provider, key)
	if err != nil {
		t.Fatalf("Failed to set API key: %v", err)
	}

	// Test GetAPIKey
	retrievedKey, err := GetAPIKey(provider)
	if err != nil {
		t.Fatalf("Failed to get API key: %v", err)
	}
	if retrievedKey != key {
		t.Errorf("Expected key %s, got %s", key, retrievedKey)
	}

	// Test DeleteAPIKey
	err = DeleteAPIKey(provider)
	if err != nil {
		t.Fatalf("Failed to delete API key: %v", err)
	}

	// Verify deletion
	_, err = GetAPIKey(provider)
	if err != keyring.ErrNotFound {
		t.Errorf("Expected ErrNotFound after deletion, got %v", err)
	}
}
