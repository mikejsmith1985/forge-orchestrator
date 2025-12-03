package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/mikejsmith1985/forge-orchestrator/internal/data"
	"github.com/mikejsmith1985/forge-orchestrator/internal/server"
)

//go:embed frontend/dist/*
var frontendEmbed embed.FS

func main() {
	// Initialize SQLite Database
	db, err := sql.Open("sqlite3", data.TokenLedgerPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize Schema
	if _, err := db.Exec(data.SQLiteSchema); err != nil {
		log.Fatal(err)
	}

	// Get the build output directory from the embed.FS
	distFS, err := fs.Sub(frontendEmbed, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Server
	srv := server.NewServer(db)
	router := srv.RegisterRoutes()

	// Serve the frontend (we need to wrap the router to handle API vs Static)
	// For simplicity, we can mount the frontend handler to the router if we change RegisterRoutes to return *http.ServeMux
	// Or we can just use the router as the main handler and add the file server to it.
	// Let's modify RegisterRoutes to allow adding more routes or handle it here.
	// Actually, the requirement says "Replace inline HTTP handler with server.NewServer()".
	// The server.RegisterRoutes returns a handler that handles /api/health and /ws.
	// We also need to serve the frontend.
	// Let's assume for now we serve API routes and fallback to frontend, or mount them separately.
	// Since RegisterRoutes returns http.Handler (likely a ServeMux), we can't easily add to it if it's opaque.
	// But it returns http.Handler.
	// Let's check server.go/routes.go again. It uses http.NewServeMux().

	// We can create a root mux here.
	rootMux := http.NewServeMux()

	// Mount API routes
	rootMux.Handle("/api/", router)
	rootMux.Handle("/ws", router)

	// Serve frontend for everything else
	rootMux.Handle("/", http.FileServer(http.FS(distFS)))

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", rootMux); err != nil {
		log.Fatal(err)
	}
}
