package server

import "net/http"

// --- Interface Definition (Requirement) ---

// PlayerStore defines the interface for storing and retrieving player scores.
// This allows mocking for tests and flexibility in implementation.
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	// Maybe add context later: e.g., RecordWin(ctx context.Context, name string)
}

// --- PlayerServer Definition (Minimal structure assumed for tests) ---

// PlayerServer holds dependencies like the PlayerStore and handles HTTP requests.
// (This would be defined in your actual server code)
type PlayerServer struct {
	Store PlayerStore
	// Assume Start() configures this handler
	Handler http.Handler
}

// NewPlayerServer creates a server (actual implementation detail)
func NewPlayerServer(store PlayerStore) *PlayerServer {
	return &PlayerServer{}
}

// Start configures routes (actual implementation detail)
func (p *PlayerServer) Start() {
	mux := http.NewServeMux()
	//... routes ...;
	p.Handler = mux
}

// ServeHTTP makes PlayerServer usable with httptest (actual implementation detail)
// func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) { p.Handler.ServeHTTP(w,r) }
