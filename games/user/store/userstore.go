package store

import (
	"errors"
	"games/user/lib"
)

// ErrUserNotFound is returned by PlayerStore when a user is not found.
var ErrUserNotFound = errors.New("user not found")

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
