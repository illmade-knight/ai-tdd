package store

import "games/user/lib"

// --- Interface Definition (Requirement) ---
// --- PlayerStore Interface (Updated) ---

// PlayerStore defines the interface for storing and retrieving player scores.
// This allows mocking for tests and flexibility in implementation.

type PlayerStore interface {
	// GetPlayerScore now returns (score, error)
	GetPlayerScore(name string) (int, error)
	RecordWin(name string)
	GetLeague() []lib.User // New method to get sorted league data
}
