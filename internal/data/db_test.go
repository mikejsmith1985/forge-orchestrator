// Package data provides database initialization and management for the Forge Orchestrator.
// This test file verifies that the database connection and table creation work correctly.
package data

import (
	"os"
	"testing"
)

// TestDatabaseConnection verifies that we can establish a connection to the SQLite database.
// This is the most basic test - if we can't connect, nothing else will work.
func TestDatabaseConnection(t *testing.T) {
	// Create a temporary database file for testing.
	// We use a temporary file so our tests don't mess with real data.
	tempDB := "test_connection.db"

	// Clean up after ourselves when the test is done.
	defer os.Remove(tempDB)

	// Try to connect to the database.
	db, err := Connect(tempDB)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// If we got here, the connection worked!
	// Let's verify we can actually talk to the database.
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

// TestDatabaseInitializationCreatesAllTables verifies that after initialization,
// all three required tables exist: token_ledger, forge_flows, and user_secrets.
// This is the main validation for Contract 1.
func TestDatabaseInitializationCreatesAllTables(t *testing.T) {
	// Create a temporary database file for testing.
	tempDB := "test_init.db"

	// Clean up after ourselves when the test is done.
	defer os.Remove(tempDB)

	// Initialize the database - this should create all tables.
	db, err := InitializeDatabase(tempDB)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Now check that each of the three required tables exists.
	// These are the tables specified in the contract:
	// 1. token_ledger - for tracking API usage costs
	// 2. forge_flows - for storing workflow definitions
	// 3. user_secrets - for storing encrypted API keys
	requiredTables := []string{"token_ledger", "forge_flows", "user_secrets"}

	for _, tableName := range requiredTables {
		exists, err := TableExists(db, tableName)
		if err != nil {
			t.Fatalf("Error checking if table %s exists: %v", tableName, err)
		}
		if !exists {
			t.Errorf("Table %s should exist after initialization, but it doesn't", tableName)
		}
	}
}

// TestDatabaseCreatesFileIfNotExists verifies that InitializeDatabase creates
// the database file if it doesn't already exist.
func TestDatabaseCreatesFileIfNotExists(t *testing.T) {
	// Use a unique filename that definitely doesn't exist.
	tempDB := "test_new_file.db"

	// Clean up after ourselves.
	defer os.Remove(tempDB)

	// Make sure the file doesn't exist before we start.
	if _, err := os.Stat(tempDB); err == nil {
		t.Fatal("Test file already exists - test cannot proceed")
	}

	// Initialize the database.
	db, err := InitializeDatabase(tempDB)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Now the file should exist.
	if _, err := os.Stat(tempDB); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}
}

// TestDatabaseInitializationIsIdempotent verifies that calling InitializeDatabase
// multiple times doesn't cause errors or duplicate tables.
// "Idempotent" means doing something twice has the same result as doing it once.
func TestDatabaseInitializationIsIdempotent(t *testing.T) {
	tempDB := "test_idempotent.db"
	defer os.Remove(tempDB)

	// Initialize once.
	db1, err := InitializeDatabase(tempDB)
	if err != nil {
		t.Fatalf("First initialization failed: %v", err)
	}
	db1.Close()

	// Initialize again - this should not fail.
	db2, err := InitializeDatabase(tempDB)
	if err != nil {
		t.Fatalf("Second initialization failed: %v", err)
	}
	defer db2.Close()

	// Tables should still exist and work correctly.
	requiredTables := []string{"token_ledger", "forge_flows", "user_secrets"}
	for _, tableName := range requiredTables {
		exists, err := TableExists(db2, tableName)
		if err != nil {
			t.Fatalf("Error checking table %s after re-init: %v", tableName, err)
		}
		if !exists {
			t.Errorf("Table %s missing after re-initialization", tableName)
		}
	}
}
