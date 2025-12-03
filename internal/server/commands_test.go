package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Initialize schema
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS command_cards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		command TEXT NOT NULL,
		description TEXT
	);
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	return db
}

func TestHandleCreateCommand(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	server := NewServer(db)
	handler := server.RegisterRoutes()

	cmd := CommandCard{
		Name:        "Test Command",
		Command:     "echo 'hello'",
		Description: "A test command",
	}
	body, _ := json.Marshal(cmd)

	req, _ := http.NewRequest("POST", "/api/commands", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var created CommandCard
	json.NewDecoder(rr.Body).Decode(&created)

	if created.ID == 0 {
		t.Errorf("Expected ID to be set")
	}
	if created.Name != cmd.Name {
		t.Errorf("Expected name %v, got %v", cmd.Name, created.Name)
	}
}

func TestHandleGetCommands(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Seed data
	_, err := db.Exec("INSERT INTO command_cards (name, command, description) VALUES (?, ?, ?)", "Cmd 1", "ls -la", "List files")
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}

	server := NewServer(db)
	handler := server.RegisterRoutes()

	req, _ := http.NewRequest("GET", "/api/commands", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var commands []CommandCard
	json.NewDecoder(rr.Body).Decode(&commands)

	if len(commands) != 1 {
		t.Errorf("Expected 1 command, got %v", len(commands))
	}
	if commands[0].Name != "Cmd 1" {
		t.Errorf("Expected name 'Cmd 1', got %v", commands[0].Name)
	}
}

func TestHandleDeleteCommand(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Seed data
	res, _ := db.Exec("INSERT INTO command_cards (name, command, description) VALUES (?, ?, ?)", "Cmd 1", "ls -la", "List files")
	id, _ := res.LastInsertId()

	server := NewServer(db)
	handler := server.RegisterRoutes()

	req, _ := http.NewRequest("DELETE", "/api/commands/"+strconv.Itoa(int(id)), nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Verify deletion
	var count int
	db.QueryRow("SELECT COUNT(*) FROM command_cards WHERE id = ?", id).Scan(&count)
	if count != 0 {
		t.Errorf("Expected command to be deleted")
	}
}
