package server

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// CommandCard represents a reusable terminal command.
type CommandCard struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	Description string `json:"description"`
}

// handleGetCommands returns all command cards.
// Educational Comment: We use a simple SELECT query to retrieve all rows.
// In a production app with many users, we'd likely need pagination or filtering here.
func (s *Server) handleGetCommands(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query("SELECT id, name, command, description FROM command_cards ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Failed to query commands: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var commands []CommandCard
	for rows.Next() {
		var c CommandCard
		if err := rows.Scan(&c.ID, &c.Name, &c.Command, &c.Description); err != nil {
			http.Error(w, "Failed to scan command: "+err.Error(), http.StatusInternalServerError)
			return
		}
		commands = append(commands, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commands)
}

// handleCreateCommand adds a new command card.
// Educational Comment: We decode the JSON body into a struct, then execute an INSERT statement.
// We return the ID of the newly created row so the frontend can update its state immediately.
func (s *Server) handleCreateCommand(w http.ResponseWriter, r *http.Request) {
	var c CommandCard
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if c.Name == "" || c.Command == "" {
		http.Error(w, "Name and Command are required", http.StatusBadRequest)
		return
	}

	res, err := s.db.Exec("INSERT INTO command_cards (name, command, description) VALUES (?, ?, ?)", c.Name, c.Command, c.Description)
	if err != nil {
		http.Error(w, "Failed to insert command: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	c.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

// handleDeleteCommand removes a command card by ID.
// Educational Comment: We parse the ID from the URL path (e.g., /api/commands/123).
// Then we execute a DELETE statement.
func (s *Server) handleDeleteCommand(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = s.db.Exec("DELETE FROM command_cards WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete command: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
