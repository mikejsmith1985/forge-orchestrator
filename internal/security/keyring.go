package security

import (
	"github.com/zalando/go-keyring"
)

const serviceName = "forge-orchestrator"

// SetAPIKey stores an API key for a specific provider in the OS keyring.
// Educational Comment: We use the `zalando/go-keyring` library which abstracts
// the underlying OS-specific keyring implementations (Keychain on macOS,
// Credential Manager on Windows, Secret Service/KWallet on Linux).
// The service name acts as a namespace for our application's keys.
func SetAPIKey(provider, key string) error {
	return keyring.Set(serviceName, provider, key)
}

// GetAPIKey retrieves an API key for a specific provider from the OS keyring.
// Educational Comment: If the item is not found, the library returns `keyring.ErrNotFound`.
// We return the raw error so the caller can handle it (e.g., by checking for ErrNotFound).
func GetAPIKey(provider string) (string, error) {
	return keyring.Get(serviceName, provider)
}

// DeleteAPIKey removes an API key for a specific provider from the OS keyring.
func DeleteAPIKey(provider string) error {
	return keyring.Delete(serviceName, provider)
}
