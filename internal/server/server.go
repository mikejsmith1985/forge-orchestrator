package server

import (
	"database/sql"

	"github.com/mikejsmith1985/forge-orchestrator/internal/llm"
)

type Server struct {
	db      *sql.DB
	gateway *llm.Gateway
}

func NewServer(db *sql.DB) *Server {
	return &Server{
		db:      db,
		gateway: llm.NewGateway(),
	}
}
