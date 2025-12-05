package server

import (
	"database/sql"

	"github.com/mikejsmith1985/forge-orchestrator/internal/llm"
)

type Server struct {
	db      *sql.DB
	gateway *llm.Gateway
	hub     *Hub
}

func NewServer(db *sql.DB) *Server {
	hub := NewHub()
	go hub.Run()
	return &Server{
		db:      db,
		gateway: llm.NewGateway(),
		hub:     hub,
	}
}
