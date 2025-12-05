package main

import (
	"database/sql"
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/mikejsmith1985/forge-orchestrator/internal/data"
	"github.com/mikejsmith1985/forge-orchestrator/internal/server"
	forgetls "github.com/mikejsmith1985/forge-orchestrator/internal/tls"
)

//go:embed frontend/dist/*
var frontendEmbed embed.FS

func main() {
	// Parse command line flags
	devTLS := flag.Bool("dev-tls", false, "Generate self-signed certificate for development")
	flag.Parse()

	// Initialize CORS configuration
	server.InitCORS()

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

	// Wrap the entire mux with CORS middleware
	handler := server.CORSMiddleware(mux)

	addr := ":8080"

	// Check for TLS configuration
	tlsCert := os.Getenv("FORGE_TLS_CERT")
	tlsKey := os.Getenv("FORGE_TLS_KEY")

	if tlsCert != "" && tlsKey != "" {
		// Production TLS with provided certificates
		log.Printf("🔒 Starting HTTPS server on %s", addr)
		if err := http.ListenAndServeTLS(addr, tlsCert, tlsKey, handler); err != nil {
			log.Fatal(err)
		}
	} else if *devTLS {
		// Development TLS with self-signed certificate
		log.Println("⚠️  Generating self-signed certificate for development")
		log.Println("⚠️  This is NOT suitable for production use!")

		certPEM, keyPEM, err := forgetls.GenerateSelfSignedCert()
		if err != nil {
			log.Fatalf("Failed to generate self-signed certificate: %v", err)
		}

		// Also save to .forge/certs for inspection
		forgetls.GenerateAndSaveCert()

		tlsConfig, err := forgetls.LoadTLSConfig(certPEM, keyPEM)
		if err != nil {
			log.Fatalf("Failed to load TLS config: %v", err)
		}

		server := &http.Server{
			Addr:      addr,
			Handler:   handler,
			TLSConfig: tlsConfig,
		}

		log.Printf("🔒 Starting HTTPS server on %s (self-signed)", addr)
		// Empty strings since we're using TLSConfig directly
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatal(err)
		}
	} else {
		// HTTP mode (no TLS)
		log.Println("⚠️  Running HTTP server (no TLS) - not recommended for production")
		log.Printf("Starting server on %s...", addr)
		if err := http.ListenAndServe(addr, handler); err != nil {
			log.Fatal(err)
		}
	}
}
