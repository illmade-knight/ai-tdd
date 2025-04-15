package server // Or your service package name + _test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Mock Implementation (For Testing) ---

// SpyPlayerStore is a mock implementation of PlayerStore for testing.
// It allows stubbing scores and spying on calls like RecordWin.
type SpyPlayerStore struct {
	// scores stores the scores for players, acting as our mock DB table.
	scores map[string]int
	// recordWinCalls stores the names of players for whom RecordWin was called.
	recordWinCalls []string
	// t is the testing context, useful for signaling errors from the mock if needed.
	t *testing.T
}

// NewSpyPlayerStore initializes a new SpyPlayerStore.
func NewSpyPlayerStore(t *testing.T) *SpyPlayerStore {
	return &SpyPlayerStore{
		scores:         make(map[string]int),
		recordWinCalls: []string{},
		t:              t,
	}
}

// GetPlayerScore retrieves the score for a player from the mock store.
// If the player is not found, it returns 0 as per requirements.
func (s *SpyPlayerStore) GetPlayerScore(name string) int {
	score, ok := s.scores[name]
	if !ok {
		// As per requirement: user initialized with 0 score if not found
		return 0
	}
	return score
}

// RecordWin simulates recording a win by incrementing the score in the mock
// store and recording the name of the player for spying purposes.
func (s *SpyPlayerStore) RecordWin(name string) {
	s.recordWinCalls = append(s.recordWinCalls, name)
	// Increment the score in the mock store as well for subsequent GETs
	s.scores[name]++
}

// Helper for tests to check RecordWin calls
func (s *SpyPlayerStore) AssertRecordWinCalledWith(expectedName string) {
	s.t.Helper()
	found := false
	for _, name := range s.recordWinCalls {
		if name == expectedName {
			found = true
			break
		}
	}
	if !found {
		s.t.Errorf("expected RecordWin to be called with '%s', but calls were %v", expectedName, s.recordWinCalls)
	}
}

// Helper to preload scores for testing GET
func (s *SpyPlayerStore) StubScore(name string, score int) {
	s.scores[name] = score
}

// --- Tests ---

// Helper function to create a PlayerServer instance for testing
// This simulates creating the server and running its setup logic (like Start).
func setupTestServer(t *testing.T) (*PlayerServer, *SpyPlayerStore) {
	t.Helper()
	store := NewSpyPlayerStore(t)

	// --- Simulate what PlayerServer.Start() would do: setup routes ---
	// In real code, this mux setup would be inside PlayerServer.Start() or a related method.
	mux := http.NewServeMux()
	// Replace handler functions with placeholder names for now
	mux.HandleFunc("GET /user/{name}/score", func(w http.ResponseWriter, r *http.Request) {
		playerName := r.PathValue("name") // Requires Go 1.22+ for PathValue
		// --- Placeholder for actual GET handler logic ---
		score := store.GetPlayerScore(playerName)
		fmt.Fprint(w, score) // Write score as string
		// Actual handler would likely involve JSON marshaling etc.
	})
	mux.HandleFunc("PUT /user/{name}/score", func(w http.ResponseWriter, r *http.Request) {
		playerName := r.PathValue("name") // Requires Go 1.22+ for PathValue
		// --- Placeholder for actual PUT handler logic ---
		store.RecordWin(playerName)
		w.WriteHeader(http.StatusAccepted) // Use Accepted for actions
	})
	// --- End Route Setup Simulation ---

	server := &PlayerServer{
		Store:   store,
		Handler: mux, // Assign the configured mux to the server's handler
	}

	// We don't call a blocking server.Start() here, we use the configured Handler directly.
	return server, store
}

func TestPlayerServer_GETScore(t *testing.T) {
	server, store := setupTestServer(t)

	// Define test cases for GET
	tests := []struct {
		name           string
		playerName     string
		initialScore   int // Score to stub in the store before the test
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Get score for existing player",
			playerName:     "Alice",
			initialScore:   5,
			expectedStatus: http.StatusOK,
			expectedBody:   "5",
		},
		{
			name:           "Get score for another existing player",
			playerName:     "Bob",
			initialScore:   10,
			expectedStatus: http.StatusOK,
			expectedBody:   "10",
		},
		{
			name:           "Get score for non-existent player",
			playerName:     "Charlie",
			initialScore:   0, // Store starts empty, effectively 0
			expectedStatus: http.StatusOK,
			expectedBody:   "0", // Should default to 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Stub the score if needed for this test case
			if tt.initialScore > 0 { // Only stub if we need a non-zero score
				store.StubScore(tt.playerName, tt.initialScore)
			}

			requestPath := fmt.Sprintf("/user/%s/score", tt.playerName)
			request, _ := http.NewRequest(http.MethodGet, requestPath, nil)
			response := httptest.NewRecorder()

			// Serve the request using the server's configured handler
			server.Handler.ServeHTTP(response, request)

			// Assert status code
			if response.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", response.Code, tt.expectedStatus)
			}

			// Assert response body
			if response.Body.String() != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %q want %q", response.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestPlayerServer_PUTScore(t *testing.T) {
	server, store := setupTestServer(t)

	playerName := "Alice"
	requestPath := fmt.Sprintf("/user/%s/score", playerName)

	request, _ := http.NewRequest(http.MethodPut, requestPath, nil) // No body needed for simple win record
	response := httptest.NewRecorder()

	// Serve the request
	server.Handler.ServeHTTP(response, request)

	// Assert status code (using Accepted for actions is common)
	expectedStatus := http.StatusAccepted
	if response.Code != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", response.Code, expectedStatus)
	}

	// Assert that the store's RecordWin method was called correctly (Spying)
	store.AssertRecordWinCalledWith(playerName)

	// Optional: Verify the score was actually incremented for subsequent GETs
	// (This tests the mock's internal logic which also helps validate the flow)
	t.Run("verify score incremented", func(t *testing.T) {
		actualScore := store.GetPlayerScore(playerName)
		expectedScore := 1 // Since it started at 0 and PUT increments
		if actualScore != expectedScore {
			t.Errorf("score after PUT was not incremented correctly in mock: got %d want %d", actualScore, expectedScore)
		}
	})

	// Test idempotency? If PUT means "record win", calling twice should increment twice.
	t.Run("calling PUT twice", func(t *testing.T) {
		// Reset calls for this subtest if necessary or adjust assertion logic
		// store.recordWinCalls = []string{} // Example reset if needed

		// Call PUT again
		request2, _ := http.NewRequest(http.MethodPut, requestPath, nil)
		response2 := httptest.NewRecorder()
		server.Handler.ServeHTTP(response2, request2)

		// Check score again
		actualScore := store.GetPlayerScore(playerName)
		expectedScore := 2 // Should be 2 now
		if actualScore != expectedScore {
			t.Errorf("score after second PUT was not incremented correctly: got %d want %d", actualScore, expectedScore)
		}

		// Check call count if needed
		if len(store.recordWinCalls) != 2 {
			t.Errorf("expected RecordWin to be called twice, got %d calls", len(store.recordWinCalls))
		}
	})
}
