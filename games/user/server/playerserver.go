package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"games/user/store"
	"log"
	"net/http"
)

// --- PlayerServer Definition (Minimal structure assumed for tests) ---

// PlayerServer holds dependencies like the PlayerStore and handles HTTP requests.
// (This would be defined in your actual server code)
type PlayerServer struct {
	store store.PlayerStore
	http.Handler
}

// NewPlayerServer creates a server (actual implementation detail)
func NewPlayerServer(store store.PlayerStore) *PlayerServer {

	ps := &PlayerServer{
		store: store,
	}
	router := http.NewServeMux()
	// Register handlers with path patterns for Go 1.22+
	router.HandleFunc("GET /players/{name}", ps.getPlayerHandler)
	router.HandleFunc("GET /players/{name}/wins", ps.getPlayerScoreHandler)
	router.HandleFunc("PUT /players/{name}/wins", ps.recordWinHandler)
	router.HandleFunc("POST /players/{name}/wins", ps.recordWinHandler)
	router.HandleFunc("GET /league", ps.getLeagueHandler)

	ps.Handler = router

	return ps
}

func (ps *PlayerServer) Close() {
	log.Printf("cleanup server")
}

// getPlayerScoreHandler handles GET /players/{name}
func (ps *PlayerServer) getPlayerScoreHandler(w http.ResponseWriter, r *http.Request) {
	// Common mistakes for HTTP handlers:
	// 1. Not setting Content-Type header correctly (e.g., application/json).
	// 2. Incorrect HTTP status codes for different outcomes (200, 404, 400, 500).
	// 3. Poor error handling: Leaking internal error details to the client or returning generic errors.
	// 4. Not validating path parameters or query parameters (PathValue handles some aspects).
	// 5. Forgetting to close request body if it's read (not applicable for GET here).

	// Get player name from path using r.PathValue (Go 1.22+)
	name := r.PathValue("name")

	if name == "" {
		// This case should ideally be prevented by robust routing or specific checks if needed,
		// as PathValue usually implies a non-empty segment match.
		// However, an explicit check can be a safeguard.
		http.Error(w, "Player name extracted from path is empty", http.StatusBadRequest)
		return
	}

	user, err := ps.store.GetPlayerScore(name)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			http.Error(w, fmt.Sprintf("Player %s not found", name), http.StatusNotFound)
		} else {
			fmt.Printf("Internal server error fetching score for %s: %v\n", name, err) // Log internal error
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		fmt.Printf("Error encoding player score response for %s: %v\n", name, err) // Log internal error
	}
}

// getPlayerScoreHandler handles GET /players/{name}
func (ps *PlayerServer) getPlayerHandler(w http.ResponseWriter, r *http.Request) {
	// Get player name from path using r.PathValue (Go 1.22+)
	name := r.PathValue("name")

	if name == "" {
		// This case should ideally be prevented by robust routing or specific checks if needed,
		// as PathValue usually implies a non-empty segment match.
		// However, an explicit check can be a safeguard.
		http.Error(w, "Player name extracted from path is empty", http.StatusBadRequest)
		return
	}

	user, err := ps.store.GetPlayerScore(name)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			http.Error(w, fmt.Sprintf("Player %s not found", name), http.StatusNotFound)
		} else {
			fmt.Printf("Internal server error fetching score for %s: %v\n", name, err) // Log internal error
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		fmt.Printf("Error encoding player score response for %s: %v\n", name, err) // Log internal error
	}
}

// recordWinHandler handles POST /players/{name}/wins
func (ps *PlayerServer) recordWinHandler(w http.ResponseWriter, r *http.Request) {
	// Get player name from path using r.PathValue (Go 1.22+)
	name := r.PathValue("name")

	if name == "" {
		http.Error(w, "Player name extracted from path is empty", http.StatusBadRequest)
		return
	}

	updatedUser, err := ps.store.RecordWin(name)
	if err != nil {
		fmt.Printf("Internal server error recording win for %s: %v\n", name, err) // Log internal error
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("updated user")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(updatedUser); err != nil {
		fmt.Printf("Error encoding record win response for %s: %v\n", name, err) // Log internal error
	}
}

// getLeagueHandler handles GET /league
func (ps *PlayerServer) getLeagueHandler(w http.ResponseWriter, r *http.Request) {
	league, err := ps.store.GetLeague()
	if err != nil {
		fmt.Printf("Internal server error fetching league: %v\n", err) // Log internal error
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(league); err != nil {
		fmt.Printf("Error encoding league response: %v\n", err) // Log internal error
	}
}
