package server

import (
	"encoding/json"
	"errors"
	"games/user/lib"
	"games/user/store"
	"net/http"
	"sort"
	"sync"
)

// ErrUserNotFound is returned by PlayerStore when a user is not found.
var ErrUserNotFound = errors.New("user not found")

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
		if errors.Is(err, ErrUserNotFound) {
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

// --- In-Memory Implementation ---

// InMemoryPlayerStore provides an in-memory implementation of PlayerStore.
// It uses a map to store scores and a mutex to handle concurrent access safely.
type InMemoryPlayerStore struct {
	// Mutex to protect concurrent reads/writes to the scores map
	mu sync.RWMutex
	// scores holds player names mapped to their scores.
	scores map[string]int
}

// NewInMemoryPlayerStore initializes and returns a new InMemoryPlayerStore.
// It creates the necessary map for storing scores.
func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		// Initialize the map
		scores: make(map[string]int),
		// The mu (RWMutex) is ready to use with its zero value
	}
}

// GetPlayerScore retrieves the score for a given player name.
// It returns the score and nil error if the player exists,
// otherwise returns 0 and ErrUserNotFound.
// Uses a read lock for concurrent safety.
func (i *InMemoryPlayerStore) GetPlayerScore(name string) (int, error) {
	// Acquire a read lock - allows multiple readers simultaneously
	i.mu.RLock()
	// Ensure the lock is released even if something panics (though unlikely here)
	defer i.mu.RUnlock()

	// Look up the score in the map
	score, ok := i.scores[name]
	if !ok {
		// Player not found
		return 0, ErrUserNotFound
	}
	// Player found
	return score, nil
}

// RecordWin increments the score for a given player name.
// If the player does not exist, they are created with a score of 1.
// Uses a write lock to ensure exclusive access during modification.
func (i *InMemoryPlayerStore) RecordWin(name string) {
	// Acquire a write lock - ensures exclusive access
	i.mu.Lock()
	// Ensure the lock is released
	defer i.mu.Unlock()

	// Get the current score (defaults to 0 if player doesn't exist in map yet)
	currentScore := i.scores[name]
	// Increment the score and update the map
	i.scores[name] = currentScore + 1
}

// Note: This implementation does not persist data. Scores are lost when the application stops.

// GetLeague (Implementation added here for completeness, though test will fail before calling it)
// In real scenario, this would be in the actual InMemoryPlayerStore implementation file.
func (i *InMemoryPlayerStore) GetLeague() []lib.User {
	i.mu.RLock()
	defer i.mu.RUnlock()

	league := make([]lib.User, 0, len(i.scores))
	for name, score := range i.scores {
		league = append(league, lib.User{Name: name, Score: score})
	}

	// Sort by score descending
	sort.Slice(league, func(a, b int) bool {
		return league[a].Score > league[b].Score
	})

	return league
}
