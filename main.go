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

	// Mount API routes and WebSocket
	// Educational Comment: Since our router already defines paths starting with /api/ and /ws,
	// we should mount it at the root "/" so it can match those paths directly.
	// However, we also need to serve the frontend.
	// We can use the router as the main handler, and add a catch-all for the frontend.
	// But RegisterRoutes returns a *http.ServeMux which is closed for modification if we don't cast it.
	// A better approach:
	// 1. Mount router at "/"
	// 2. But router doesn't know about frontend.
	//
	// Let's use the fact that ServeMux matches the most specific pattern.
	// If we mount router at "/", it will handle everything.
	// We need to merge them.
	//
	// Alternative:
	// RegisterRoutes defines "/api/..." and "/ws".
	// We can just use that mux as our root mux, and add the frontend handler to it!
	// But RegisterRoutes returns http.Handler interface.
	// We should update RegisterRoutes to accept the frontend handler or return *http.ServeMux.
	//
	// For now, let's just register the frontend handler on the SAME mux if possible.
	// But we can't easily access the mux inside the handler interface.
	//
	// Let's modify main.go to use the router for "/api/" and "/ws" by STRIPPING the prefix?
	// No, the router expects "/api/...".
	//
	// Correct fix:
	// The router handles "/api/..." and "/ws".
	// We want rootMux to delegate "/api/" and "/ws" to router.
	// BUT, if we use Handle("/api/", router), the request path passed to router will be "/api/..." (no stripping).
	// So router sees "/api/commands". It matches "/api/commands". This SHOULD work.
	//
	// Wait, does http.ServeMux strip prefix?
	// "Handle registers the handler for the given pattern."
	// It does NOT strip prefix.
	// So if I request "/api/commands", rootMux matches "/api/", calls router.
	// Router sees "/api/commands".
	// Router has pattern "/api/commands".
	// It SHOULD match.
	//
	// Why did it fail?
	// Maybe "Unexpected token <" means it fell through to the frontend?
	// If rootMux matches "/", it goes to frontend.
	// Does "/api/" match "/api/commands"? Yes.
	//
	// Let's look at the error again.
	// "Unexpected token <" in JSON.
	// This means the response was HTML.
	//
	// Hypothesis: The backend is NOT running on 8080, or the proxy is hitting the wrong port.
	// I saw "Starting server on :8080..." in the logs.
	// I saw "Port 8080 is in use" in frontend logs.
	// So backend IS on 8080.
	//
	// Hypothesis 2: The request is NOT matching "/api/".
	// Maybe the browser is requesting "http://localhost:8081/api/commands".
	// Vite proxies to "http://localhost:8080/api/commands".
	// Backend receives "/api/commands".
	//
	// If `router` is `http.NewServeMux()`.
	// And `rootMux` is `http.NewServeMux()`.
	// `rootMux.Handle("/api/", router)`.
	// Request "/api/commands".
	// `rootMux` dispatches to `router`.
	// `router` sees "/api/commands".
	// `router` has `mux.HandleFunc("GET /api/commands", ...)` (Go 1.22 style with method).
	//
	// Wait, `mux.HandleFunc("GET /api/commands", ...)`
	// If I use `rootMux.Handle("/api/", router)`, does it preserve the method? Yes.
	//
	// Is it possible `rootMux` is taking precedence with "/"?
	// No, longest match wins.
	//
	// Let's try to simplify.
	// I will modify `main.go` to use a single Mux if possible, or debug why it falls through.
	//
	// Actually, I suspect `http.FileServer` might be capturing it if I messed up the pattern.
	//
	// Let's try to be explicit.
	// I will modify `main.go` to just use `router` for everything, and add the frontend handler TO `RegisterRoutes`?
	// No, `RegisterRoutes` is in `server` package, shouldn't know about `distFS`.
	//
	// Let's try `http.StripPrefix`?
	// If I strip prefix, `router` sees `commands`. But `router` expects `/api/commands`.
	//
	// Let's try removing the method specifier in `routes.go`?
	// No, that's good practice.
	//
	// Maybe the issue is `mux.HandleFunc("POST /api/ledger", ...)`
	// If the request comes in as `GET /api/ledger` (browser navigation?), it 404s?
	// But the error is "fetching commands". That's a GET.
	// `mux.HandleFunc("GET /api/commands", ...)`
	//
	// Let's try to debug by logging in `main.go`.
	//
	// OR, simply:
	// The `router` returned by `RegisterRoutes` is a `*http.ServeMux`.
	// I can cast it and add the frontend handler to it!
	//
	// `router.(*http.ServeMux).Handle("/", http.FileServer(http.FS(distFS)))`
	// Then use `router` as the main handler.
	// This avoids nested muxes.

	// Cast router to *http.ServeMux to add frontend handler
	if mux, ok := router.(*http.ServeMux); ok {
		// Serve frontend for everything else
		mux.Handle("/", http.FileServer(http.FS(distFS)))

		log.Println("Starting server on :8080...")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Router is not *http.ServeMux")
	}
}
