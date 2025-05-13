package store

import (
	"games/user/lib"
	"sort"
	"sync"
)

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
func (i *InMemoryPlayerStore) RecordWin(name string) (*lib.User, error) {
	// Acquire a write lock - ensures exclusive access
	i.mu.Lock()
	// Ensure the lock is released
	defer i.mu.Unlock()

	// Get the current score (defaults to 0 if player doesn't exist in map yet)
	currentScore := i.scores[name]
	// Increment the score and update the map
	i.scores[name] = currentScore + 1
	return &lib.User{
		Name:  name,
		Score: i.scores[name],
	}, nil
}

// Note: This implementation does not persist data. Scores are lost when the application stops.

// GetLeague (Implementation added here for completeness, though test will fail before calling it)
// In real scenario, this would be in the actual InMemoryPlayerStore implementation file.
func (i *InMemoryPlayerStore) GetLeague() ([]lib.User, error) {
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

	return league, nil
}
