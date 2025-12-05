// Package data provides database initialization and management for the Forge Orchestrator.
// This file contains functions to connect to and initialize the SQLite database.
package data

import (
	"database/sql"
	"os"

	// Import the SQLite driver. The underscore means we only want the driver to register itself.
	// We don't use the package directly - it just needs to be available for database/sql.
	_ "github.com/mattn/go-sqlite3"
)

// DB is the database connection that can be used throughout the application.
// It's a pointer to a sql.DB, which is Go's standard way to talk to databases.
var DB *sql.DB

// Connect opens a connection to the SQLite database file.
// If the file doesn't exist, SQLite will create it automatically.
// Think of it like opening a notebook - if you don't have one, you get a new blank one.
func Connect(dbPath string) (*sql.DB, error) {
	// sql.Open prepares a connection to the database.
	// "sqlite3" tells Go which type of database we're connecting to.
	// dbPath is the file path where our data will be stored.
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Ping checks if we can actually talk to the database.
	// It's like saying "hello, are you there?" to make sure everything works.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// InitializeDatabase creates the database file if it doesn't exist and sets up all the tables.
// Tables are like spreadsheets in Excel - they organize our data into rows and columns.
// This function creates three important tables:
// 1. token_ledger - tracks how much we spend on AI API calls
// 2. forge_flows - stores the workflow diagrams users create
// 3. user_secrets - keeps API keys safe and encrypted
func InitializeDatabase(dbPath string) (*sql.DB, error) {
	// First, we check if the database file already exists.
	// If it doesn't exist, SQLite will create it when we connect.
	_, err := os.Stat(dbPath)
	fileExists := err == nil

	// Connect to the database (creates file if needed).
	db, err := Connect(dbPath)
	if err != nil {
		return nil, err
	}

	// If the file existed, the tables might already be there.
	// But we use "CREATE TABLE IF NOT EXISTS" in our schema,
	// so running the schema again is safe - it won't duplicate tables.
	if !fileExists {
		// Brand new database, so we definitely need to create tables.
		if _, err := db.Exec(SQLiteSchema); err != nil {
			return nil, err
		}
	} else {
		// Database file exists, but let's make sure all tables are present.
		// This handles the case where the app was stopped mid-initialization.
		if _, err := db.Exec(SQLiteSchema); err != nil {
			return nil, err
		}
	}

	// Store the connection globally so other parts of the app can use it.
	DB = db

	return db, nil
}

// TableExists checks if a specific table exists in the database.
// This is useful for testing and validation to confirm our schema was applied correctly.
// It returns true if the table exists, false if it doesn't.
func TableExists(db *sql.DB, tableName string) (bool, error) {
	// SQLite stores information about tables in a special table called sqlite_master.
	// We're asking: "Is there a table with this name?"
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name=?"

	var name string
	err := db.QueryRow(query, tableName).Scan(&name)

	if err == sql.ErrNoRows {
		// No rows means the table doesn't exist. That's not an error, just a "no".
		return false, nil
	}
	if err != nil {
		// Something else went wrong (like database connection issues).
		return false, err
	}

	// If we got here, we found the table.
	return true, nil
}
