package server

import (
	"database/sql"
)

type Server struct {
	db *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{
		db: db,
	}
}
