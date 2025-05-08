package store

import (
	"errors"
	"games/user/lib"
	"path/filepath"
	"reflect"
	"testing"
)

// --- Helper to create store with temp file ---

// createTempParquetStore creates a ParquetPlayerStore using a temporary file.
// It returns the store instance and the path to the temp file.
// It uses t.Fatalf on errors during setup.
func createTempParquetStore(t *testing.T) (*ParquetPlayerStore, string) {
	t.Helper() // Mark as test helper

	tempDir := t.TempDir() // Create temp dir, automatically cleaned up
	tempFilePath := filepath.Join(tempDir, "test_players.parquet")

	// Ensure NewParquetPlayerStore is accessible (might need import or be in same package)
	store, err := NewParquetPlayerStore(tempFilePath)
	if err != nil {
		t.Fatalf("Failed to create ParquetPlayerStore for test: %v", err)
	}
	return store, tempFilePath
}

// --- Unit Tests ---

func TestParquetPlayerStore_GetPlayerScore(t *testing.T) {
	t.Run("returns not found error for non-existent user in empty store", func(t *testing.T) {
		store, _ := createTempParquetStore(t)

		_, err := store.GetPlayerScore("NonExistentUser")

		// Assert error is ErrUserNotFound
		if !errors.Is(err, ErrUserNotFound) {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("returns score and nil error for existing user", func(t *testing.T) {
		store, _ := createTempParquetStore(t)
		playerName := "PlayerOne"
		expectedScore := 3

		// Pre-populate store
		for i := 0; i < expectedScore; i++ {
			store.RecordWin(playerName)
		}

		// Get score
		score, err := store.GetPlayerScore(playerName)

		// Assert error is nil
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		// Assert score is correct
		if score != expectedScore {
			t.Errorf("expected score %d, got %d", expectedScore, score)
		}
	})

	t.Run("returns not found error for non-existent user in non-empty store", func(t *testing.T) {
		store, _ := createTempParquetStore(t)
		store.RecordWin("ExistingPlayer") // Add some data

		_, err := store.GetPlayerScore("NonExistentUser")

		// Assert error is ErrUserNotFound
		if !errors.Is(err, ErrUserNotFound) {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})
}

func TestParquetPlayerStore_RecordWin(t *testing.T) {
	t.Run("creates user with score 1 on first win", func(t *testing.T) {
		store, _ := createTempParquetStore(t)
		playerName := "NewPlayer"

		store.RecordWin(playerName)

		score, err := store.GetPlayerScore(playerName)
		if err != nil {
			t.Fatalf("GetPlayerScore failed after RecordWin: %v", err)
		}
		if score != 1 {
			t.Errorf("expected score 1 after first win, got %d", score)
		}
	})

	t.Run("increments score for existing user", func(t *testing.T) {
		store, _ := createTempParquetStore(t)
		playerName := "ExistingPlayer"

		store.RecordWin(playerName) // Score 1
		store.RecordWin(playerName) // Score 2
		store.RecordWin(playerName) // Score 3

		score, err := store.GetPlayerScore(playerName)
		if err != nil {
			t.Fatalf("GetPlayerScore failed after multiple RecordWin calls: %v", err)
		}
		if score != 3 {
			t.Errorf("expected score 3 after multiple wins, got %d", score)
		}
	})

	t.Run("handles multiple users correctly", func(t *testing.T) {
		store, _ := createTempParquetStore(t)
		playerA := "PlayerA"
		playerB := "PlayerB"

		store.RecordWin(playerA) // A: 1
		store.RecordWin(playerB) // B: 1
		store.RecordWin(playerA) // A: 2

		scoreA, errA := store.GetPlayerScore(playerA)
		scoreB, errB := store.GetPlayerScore(playerB)

		if errA != nil || errB != nil {
			t.Fatalf("GetPlayerScore failed for multiple users: errA=%v, errB=%v", errA, errB)
		}
		if scoreA != 2 {
			t.Errorf("expected score 2 for %s, got %d", playerA, scoreA)
		}
		if scoreB != 1 {
			t.Errorf("expected score 1 for %s, got %d", playerB, scoreB)
		}
	})

	// Optional: Add test to directly verify file content after RecordWin
	// This requires importing parquet-go and reading the file.
	// t.Run("writes correct data to parquet file", func(t *testing.T) { ... })
}

func TestParquetPlayerStore_GetLeague(t *testing.T) {
	t.Run("returns empty slice for empty store", func(t *testing.T) {
		store, _ := createTempParquetStore(t)
		league := store.GetLeague()

		if len(league) != 0 {
			t.Errorf("expected empty league, got %d items", len(league))
		}
	})

	t.Run("returns league sorted by score descending", func(t *testing.T) {
		store, _ := createTempParquetStore(t)

		// Add players out of order
		store.RecordWin("Charlie") // Score 1
		store.RecordWin("Alice")   // Score 1
		store.RecordWin("Alice")   // Score 2
		store.RecordWin("Bob")     // Score 1
		store.RecordWin("Bob")     // Score 2
		store.RecordWin("Bob")     // Score 3

		// Expected league order: Bob (3), Alice (2), Charlie (1)
		expectedLeague := []lib.User{
			{Name: "Bob", Score: 3},
			{Name: "Alice", Score: 2},
			{Name: "Charlie", Score: 1},
		}

		// Get league from store
		gotLeague := store.GetLeague()

		// Assert slices are equal (order matters here due to sorting requirement)
		if !reflect.DeepEqual(gotLeague, expectedLeague) {
			t.Errorf("GetLeague returned incorrect or unsorted data:\ngot:  %+v\nwant: %+v", gotLeague, expectedLeague)
		}
	})

	t.Run("returns correct league after multiple updates", func(t *testing.T) {
		store, _ := createTempParquetStore(t)

		store.RecordWin("Zoe")    // 1
		store.RecordWin("Xavier") // 1
		store.RecordWin("Zoe")    // 2
		store.RecordWin("Yannis") // 1
		store.RecordWin("Zoe")    // 3
		store.RecordWin("Xavier") // 2

		expectedLeague := []lib.User{
			{Name: "Zoe", Score: 3},
			{Name: "Xavier", Score: 2},
			{Name: "Yannis", Score: 1},
		}

		gotLeague := store.GetLeague()

		if !reflect.DeepEqual(gotLeague, expectedLeague) {
			t.Errorf("GetLeague returned incorrect or unsorted data after updates:\ngot:  %+v\nwant: %+v", gotLeague, expectedLeague)
		}
	})
}
