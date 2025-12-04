package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

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

	// Cast to *http.ServeMux to add SPA handler
	mux, ok := router.(*http.ServeMux)
	if !ok {
		log.Fatal("Router is not *http.ServeMux")
	}

	// SPA Handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Prepare path for fs.Open (no leading slash)
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "."
		}

		// Try to open the file to check if it exists
		f, err := distFS.Open(path)
		if err != nil {
			// If file not found, assume SPA route and serve index.html
			r.URL.Path = "/"
		} else {
			defer f.Close()
		}

		http.FileServer(http.FS(distFS)).ServeHTTP(w, r)
	})

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
