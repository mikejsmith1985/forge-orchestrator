package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestHealthHandler(t *testing.T) {
	// Create a new server instance (mocking DB as nil for now since health check doesn't use it)
	s := NewServer(nil)

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := s.RegisterRoutes()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"status": "ok"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestWebSocketHandler(t *testing.T) {
	s := NewServer(nil)
	server := httptest.NewServer(s.RegisterRoutes())
	defer server.Close()

	// Convert http URL to ws URL
	u := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer ws.Close()

	// Send a message
	message := []byte("hello")
	err = ws.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		t.Fatalf("write: %v", err)
	}

	// Read the message back (echo)
	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	if string(p) != string(message) {
		t.Errorf("handler returned unexpected message: got %v want %v",
			string(p), string(message))
	}
}
