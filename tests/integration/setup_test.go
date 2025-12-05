package integration

import (
	"database/sql"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mikejsmith1985/forge-orchestrator/internal/data"
	"github.com/mikejsmith1985/forge-orchestrator/internal/server"
	_ "github.com/mattn/go-sqlite3"
)

var testServer *httptest.Server
var testDB *sql.DB

// TestMain sets up the integration test environment with a real server
// and in-memory SQLite database. This allows testing the full stack
// without mocking API responses.
func TestMain(m *testing.M) {
	// Setup in-memory SQLite database
	var err error
	testDB, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		fmt.Printf("Failed to open in-memory database: %v\n", err)
		os.Exit(1)
	}

	// Initialize schema
	_, err = testDB.Exec(data.SQLiteSchema)
	if err != nil {
		fmt.Printf("Failed to initialize schema: %v\n", err)
		os.Exit(1)
	}

	// Create server with real database
	srv := server.NewServer(testDB)
	testServer = httptest.NewServer(srv.RegisterRoutesWithCORS())

	// Run tests
	code := m.Run()

	// Teardown
	testServer.Close()
	testDB.Close()

	os.Exit(code)
}

// GetServerURL returns the URL of the test server for Playwright tests
func GetServerURL() string {
	return testServer.URL
}

// GetDB returns the test database for direct data verification
func GetDB() *sql.DB {
	return testDB
}

// ResetDB clears all data from the test database between tests
func ResetDB(t *testing.T) {
	tables := []string{"command_cards", "forge_flows", "token_ledger", "optimization_suggestions", "user_secrets"}
	for _, table := range tables {
		_, err := testDB.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			t.Fatalf("Failed to reset table %s: %v", table, err)
		}
	}
}
