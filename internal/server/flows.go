package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mikejsmith1985/forge-orchestrator/internal/flows"
)

// handleGetFlows retrieves all flows.
func (s *Server) handleGetFlows(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, name, data, status, created_at FROM forge_flows ORDER BY created_at DESC`
	rows, err := s.db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []flows.Flow
	for rows.Next() {
		var f flows.Flow
		if err := rows.Scan(&f.ID, &f.Name, &f.Data, &f.Status, &f.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, f)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleGetFlow retrieves a single flow by ID.
func (s *Server) handleGetFlow(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `SELECT id, name, data, status, created_at FROM forge_flows WHERE id = ?`
	var f flows.Flow
	err = s.db.QueryRow(query, id).Scan(&f.ID, &f.Name, &f.Data, &f.Status, &f.CreatedAt)
	if err != nil {
		http.Error(w, "Flow not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(f)
}

// handleCreateFlow creates a new flow.
func (s *Server) handleCreateFlow(w http.ResponseWriter, r *http.Request) {
	var f flows.Flow
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO forge_flows (name, data, status) VALUES (?, ?, ?)`
	res, err := s.db.Exec(query, f.Name, f.Data, f.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	f.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(f)
}

// handleUpdateFlow updates an existing flow.
func (s *Server) handleUpdateFlow(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var f flows.Flow
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE forge_flows SET name = ?, data = ?, status = ? WHERE id = ?`
	_, err = s.db.Exec(query, f.Name, f.Data, f.Status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// handleDeleteFlow deletes a flow.
func (s *Server) handleDeleteFlow(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM forge_flows WHERE id = ?`
	_, err = s.db.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleExecuteFlow triggers the execution of a flow.
func (s *Server) handleExecuteFlow(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Create file signaler for fallback
	fileSignaler, _ := flows.NewFileSignaler()

	// Execute the flow with Hub integration for real-time broadcasts
	if err := flows.ExecuteFlowWithHub(id, s.db, s.gateway, nil, fileSignaler, s.hub); err != nil {
		http.Error(w, "Flow execution failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "completed"}`))
}
