package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	_ "modernc.org/sqlite"

	"github.com/mikejsmith1985/forge-orchestrator/internal/config"
	"github.com/mikejsmith1985/forge-orchestrator/internal/data"
	"github.com/mikejsmith1985/forge-orchestrator/internal/server"
	forgetls "github.com/mikejsmith1985/forge-orchestrator/internal/tls"
	"github.com/mikejsmith1985/forge-orchestrator/internal/updater"
)

//go:embed frontend/dist/*
var frontendEmbed embed.FS

// Preferred ports to try, in order
var preferredPorts = []int{8080, 8333, 9000, 3000, 3333}

func main() {
	// Parse command line flags
	devTLS := flag.Bool("dev-tls", false, "Generate self-signed certificate for development")
	noBrowser := flag.Bool("no-browser", false, "Don't open browser on startup")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Warning: Failed to load config, using defaults: %v", err)
		cfg = config.DefaultConfig()
	}

	// Initialize CORS configuration
	server.InitCORS()

	// Get database path from config
	dbPath, err := config.GetDatabasePath()
	if err != nil {
		log.Printf("Warning: Failed to get data directory, using local path: %v", err)
		dbPath = "forge_ledger.db"
	}

	// Initialize SQLite Database
	db, err := sql.Open("sqlite", dbPath)
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

	// Cast to *http.ServeMux to add handlers
	mux, ok := router.(*http.ServeMux)
	if !ok {
		log.Fatal("Router is not *http.ServeMux")
	}

	// Add update API endpoints
	mux.HandleFunc("/api/version", handleVersion)
	mux.HandleFunc("/api/update/check", handleUpdateCheck)
	mux.HandleFunc("/api/update/apply", handleUpdateApply)
	mux.HandleFunc("/api/update/versions", handleListVersions)

	// Add config API endpoints
	mux.HandleFunc("/api/config", handleConfig)
	mux.HandleFunc("/api/wsl/detect", handleWSLDetect)

	// Add shutdown endpoint
	mux.HandleFunc("/api/shutdown", handleShutdown)

	// SPA Handler (must be last)
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

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		log.Println("\n👋 Shutting down Forge Orchestrator...")
		os.Exit(0)
	}()

	// Check for TLS configuration
	tlsCert := os.Getenv("FORGE_TLS_CERT")
	tlsKey := os.Getenv("FORGE_TLS_KEY")

	if tlsCert != "" && tlsKey != "" {
		// Production TLS with provided certificates
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Printf("🔒 Starting HTTPS server on %s", addr)
		if err := http.ListenAndServeTLS(addr, tlsCert, tlsKey, handler); err != nil {
			log.Fatal(err)
		}
	} else if *devTLS {
		// Development TLS with self-signed certificate
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
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
		// HTTP mode with port fallback
		addr, listener, err := findAvailablePort(cfg.Server.Port)
		if err != nil {
			log.Fatalf("Failed to find available port: %v", err)
		}

		log.Printf("🔥 Forge Orchestrator v%s starting at http://%s", updater.GetVersion(), addr)
		log.Printf("📁 Database: %s", dbPath)

		// Auto-open browser (unless disabled)
		if cfg.Server.OpenBrowser && !*noBrowser && os.Getenv("NO_BROWSER") == "" {
			go openBrowser("http://" + addr)
		}

		if err := http.Serve(listener, handler); err != nil {
			log.Fatal(err)
		}
	}
}

// findAvailablePort tries the preferred port first, then falls back to alternatives
func findAvailablePort(preferred int) (string, net.Listener, error) {
	// Try preferred port first
	ports := []int{preferred}
	for _, p := range preferredPorts {
		if p != preferred {
			ports = append(ports, p)
		}
	}

	for _, port := range ports {
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			return addr, listener, nil
		}
		log.Printf("Port %d unavailable, trying next...", port)
	}

	// Fallback: let OS assign a random available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, fmt.Errorf("no available ports: %w", err)
	}
	addr := listener.Addr().String()
	log.Printf("Using OS-assigned port: %s", addr)
	return addr, listener, nil
}

// openBrowser opens the default browser to the given URL
func openBrowser(url string) {
	time.Sleep(500 * time.Millisecond) // Small delay to let server start
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	}
	if cmd != nil {
		_ = cmd.Start()
	}
}

// API Handlers

func handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"version": updater.GetVersion(),
	})
}

func handleUpdateCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	info, err := updater.CheckForUpdate()
	if err != nil {
		log.Printf("[Updater] Check failed: %v", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"available":      false,
			"currentVersion": updater.GetVersion(),
			"error":          err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(info)
}

func handleUpdateApply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Check for update first
	info, err := updater.CheckForUpdate()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if !info.Available {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "No update available",
		})
		return
	}

	// Download the update
	log.Printf("[Updater] Downloading %s...", info.AssetName)
	tmpPath, err := updater.DownloadUpdate(info)
	if err != nil {
		log.Printf("[Updater] Download failed: %v", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Download failed: " + err.Error(),
		})
		return
	}

	// Apply the update
	log.Printf("[Updater] Applying update...")
	if err := updater.ApplyUpdate(tmpPath); err != nil {
		log.Printf("[Updater] Apply failed: %v", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Apply failed: " + err.Error(),
		})
		return
	}

	log.Printf("[Updater] Update applied successfully! Restarting...")

	// Send success response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"newVersion": info.LatestVersion,
		"message":    "Update applied. Restarting...",
	})

	// Restart the application
	go func() {
		time.Sleep(500 * time.Millisecond)
		restartSelf()
	}()
}

func restartSelf() {
	executable, err := os.Executable()
	if err != nil {
		log.Printf("[Updater] Failed to get executable path: %v", err)
		os.Exit(1)
	}

	// On Windows, we need to start a new process and exit
	// On Unix, we can use exec to replace the current process
	if runtime.GOOS == "windows" {
		cmd := exec.Command(executable)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
		os.Exit(0)
	} else {
		syscall.Exec(executable, []string{executable}, os.Environ())
	}
}

func handleListVersions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	releases, err := updater.ListReleases(10) // Get last 10 releases
	if err != nil {
		log.Printf("[Updater] Failed to list releases: %v", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":    err.Error(),
			"releases": []interface{}{},
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"releases":       releases,
		"currentVersion": updater.GetVersion(),
	})
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		cfg, err := config.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(cfg)

	case http.MethodPost:
		var cfg config.Config
		if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := config.Save(&cfg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleWSLDetect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if runtime.GOOS != "windows" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"available": false,
			"reason":    "Not running on Windows",
		})
		return
	}

	// Get list of WSL distros
	cmd := exec.Command("wsl", "--list", "--quiet")
	output, err := cmd.Output()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"available": false,
			"reason":    "WSL not installed or not available",
		})
		return
	}

	// Parse distro names (handle UTF-16 output from wsl.exe)
	distros := []string{}
	lines := strings.Split(strings.ReplaceAll(string(output), "\x00", ""), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			distros = append(distros, line)
		}
	}

	if len(distros) == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"available": false,
			"reason":    "No WSL distributions installed",
		})
		return
	}

	// Try to get the username from the first distro
	username := ""
	if len(distros) > 0 {
		userCmd := exec.Command("wsl", "-d", distros[0], "-e", "whoami")
		userOutput, err := userCmd.Output()
		if err == nil {
			username = strings.TrimSpace(string(userOutput))
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"available":   true,
		"distros":     distros,
		"defaultUser": username,
		"defaultHome": "/home/" + username,
	})
}

func handleShutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"shutting down"}`))
	log.Println("👋 Shutdown requested from browser")
	go func() {
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}()
}
