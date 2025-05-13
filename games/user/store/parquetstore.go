package store

import (
	"errors"
	"fmt"
	"games/user/lib"
	"github.com/parquet-go/parquet-go"
	"io"
	"log"
	"os"
	"sort"
	"sync"
)

// --- Parquet Implementation ---

// ParquetPlayerStore provides a PlayerStore implementation using a Parquet file.
type ParquetPlayerStore struct {
	filePath string
	mu       sync.RWMutex // Mutex to protect file access
}

// NewParquetPlayerStore creates a new store instance linked to a file path.
func NewParquetPlayerStore(filePath string) (*ParquetPlayerStore, error) {
	store := &ParquetPlayerStore{
		filePath: filePath,
	}
	// Optional: Try an initial load to validate the file/path early
	_, err := store.loadData()
	if err != nil && !errors.Is(err, os.ErrNotExist) { // Ignore "not found" error on init
		return nil, fmt.Errorf("failed to initialize parquet store from %s: %w", filePath, err)
	}
	return store, nil
}

// loadData reads the parquet file into an in-memory map.
// IMPORTANT: Inefficient for large files - reads the whole file.
func (p *ParquetPlayerStore) loadData() (map[string]int, error) {
	// Open the file for reading
	file, err := os.Open(p.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// File doesn't exist yet, return empty map, not an error
			return make(map[string]int), nil
		}
		return nil, fmt.Errorf("error opening parquet file %s: %w", p.filePath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close parquet file %s: %v", p.filePath, err)
		}
	}(file)

	// Create a parquet reader
	reader := parquet.NewReader(file) // Using default options

	// Read all rows into a slice of User structs
	// Note: For very large files, reading row-by-row might be better,
	// but requires knowing the number of rows or handling EOF.
	// Reading all is simpler for this example.
	rows := make([]lib.User, 0, reader.NumRows()) // Pre-allocate slice if possible
	for {
		var user lib.User
		err := reader.Read(&user)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break // End of file reached
			}
			// Check for schema mismatch or other read errors
			return nil, fmt.Errorf("error reading row from parquet file %s: %w", p.filePath, err)
		}
		rows = append(rows, user)
	}

	// Convert slice to map for easier lookup
	scores := make(map[string]int, len(rows))
	for _, user := range rows {
		scores[user.Name] = user.Score
	}

	return scores, nil
}

// saveData writes the current scores map back to the parquet file.
// IMPORTANT: Inefficient - writes the whole file.
func (p *ParquetPlayerStore) saveData(scores map[string]int) error {
	// Create or truncate the file for writing
	file, err := os.Create(p.filePath) // os.Create truncates if file exists
	if err != nil {
		return fmt.Errorf("error creating/truncating parquet file %s: %w", p.filePath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close parquet file %s: %v", p.filePath, err)
		}
	}(file)

	// Convert map back to slice of User structs
	league := make([]lib.User, 0, len(scores))
	for name, score := range scores {
		league = append(league, lib.User{Name: name, Score: score})
	}

	// Create a parquet writer
	// Using default writer options. Compression (e.g., ZSTD, SNAPPY) is recommended.
	writer := parquet.NewWriter(file, parquet.SchemaOf(lib.User{})) // Infer schema from User struct

	// Write all users to the file
	for _, user := range league {
		if err := writer.Write(user); err != nil {
			// Attempt to close writer even on error
			_ = writer.Close()
			return fmt.Errorf("error writing user %s to parquet file %s: %w", user.Name, p.filePath, err)
		}
	}

	// Close the writer to flush buffers and write footers
	if err := writer.Close(); err != nil {
		return fmt.Errorf("error closing parquet writer for %s: %w", p.filePath, err)
	}

	return nil
}

// GetPlayerScore retrieves the score for a given player name from the file.
func (p *ParquetPlayerStore) GetPlayerScore(name string) (int, error) {
	p.mu.RLock() // Read lock for loading data
	scores, err := p.loadData()
	p.mu.RUnlock() // Release lock after loading

	if err != nil {
		// Log the load error, but try to continue if possible (might be transient?)
		// Or return a specific internal error? For now, log and treat as not found.
		log.Printf("WARN: Error loading data in GetPlayerScore: %v", err)
		return 0, ErrUserNotFound // Or return the actual load error? Interface doesn't specify.
	}

	score, ok := scores[name]
	if !ok {
		return 0, ErrUserNotFound
	}
	return score, nil
}

// RecordWin increments the score for a given player name in the file.
func (p *ParquetPlayerStore) RecordWin(name string) (*lib.User, error) {
	p.mu.Lock() // Write lock for load/modify/save sequence
	defer p.mu.Unlock()

	scores, err := p.loadData()
	if err != nil {
		// Log critical error - cannot reliably update score if load fails
		log.Printf("ERROR: Failed to load data in RecordWin for %s: %v. Score not updated.", name, err)
		return nil, ErrDataSource
	}

	// Increment score (creates user with score 1 if not present)
	scores[name]++

	// Save updated data back to file
	err = p.saveData(scores)
	if err != nil {
		// Log critical error - update was lost
		log.Printf("ERROR: Failed to save data in RecordWin for %s after increment: %v. Update lost.", name, err)
		// How to handle this? Rollback memory change? The interface doesn't allow returning errors.
		// For now, the in-memory map *was* changed, but the file write failed.
		// This highlights the limitations of the simple load/save approach.
	}
	log.Println("saved win")

	return &lib.User{Name: name, Score: scores[name]}, nil
}

// GetLeague reads the file and returns all players sorted by score.
func (p *ParquetPlayerStore) GetLeague() ([]lib.User, error) {
	p.mu.RLock() // Read lock for loading data
	scores, err := p.loadData()
	p.mu.RUnlock() // Release lock after loading

	if err != nil {
		log.Printf("ERROR: Failed to load data in GetLeague: %v. Returning empty league.", err)
		return []lib.User{}, ErrDataSource // Return empty slice on error
	}

	// Convert map to slice
	league := make([]lib.User, 0, len(scores))
	for name, score := range scores {
		league = append(league, lib.User{Name: name, Score: score})
	}

	// Sort by score descending
	sort.Slice(league, func(a, b int) bool {
		return league[a].Score > league[b].Score
	})

	return league, nil
}
