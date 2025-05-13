package store

import (
	"errors"
	"games/user/lib"
)

// ErrUserNotFound is returned by PlayerStore when a user is not found.
var ErrUserNotFound = errors.New("user not found")

// ErrDataSource is returned by PlayerStore when data cannot be retrieved (details will be implementation specific)
var ErrDataSource = errors.New("could not access data source")

// --- Interface Definition (Requirement) ---
// --- PlayerStore Interface (Updated) ---

// PlayerStore defines the interface for storing and retrieving player scores.
// This allows mocking for tests and flexibility in implementation.

type PlayerStore interface {
	// GetPlayerScore now returns (score, error)
	GetPlayerScore(name string) (int, error)
	RecordWin(name string) (*lib.User, error)
	GetLeague() ([]lib.User, error) // New method to get sorted league data
}
