package server

import (
	"encoding/json"
	"errors"
	"games/user/lib"
	"games/user/store"
	"net/http"
)

// --- PlayerServer Definition (Minimal structure assumed for tests) ---

// PlayerServer holds dependencies like the PlayerStore and handles HTTP requests.
// (This would be defined in your actual server code)
type PlayerServer struct {
	store store.PlayerStore
	// Assume Start() configures this handler
	Handler http.Handler
}

// NewPlayerServer creates a server (actual implementation detail)
func NewPlayerServer(store store.PlayerStore) *PlayerServer {
	return &PlayerServer{store: store}
}

// Start configures routes (actual implementation detail)
func (p *PlayerServer) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /user/{name}/score", func(w http.ResponseWriter, r *http.Request) {
		playerName := r.PathValue("name") // Go 1.22+
		score, err := p.store.GetPlayerScore(playerName)
		if errors.Is(err, store.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		user := lib.User{Name: playerName, Score: score}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	})

	mux.HandleFunc("PUT /user/{name}/score", func(w http.ResponseWriter, r *http.Request) {
		playerName := r.PathValue("name") // Go 1.22+
		p.store.RecordWin(playerName)
		w.WriteHeader(http.StatusAccepted)
	})

	p.Handler = mux
}
