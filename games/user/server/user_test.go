package server // Or your service package name + _test

import (
	"encoding/json"
	"errors"
	"fmt"
	"games/user/lib"
	"games/user/store"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
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
// GetPlayerScore (Updated) retrieves the score or returns ErrUserNotFound.
func (s *SpyPlayerStore) GetPlayerScore(name string) (int, error) {
	score, ok := s.scores[name]
	if !ok {
		// Return 0 score and the specific error
		return 0, store.ErrUserNotFound
	}
	// lib.User found, return score and nil error
	return score, nil
}

// RecordWin simulates recording a win by incrementing the score in the mock
// store and recording the name of the player for spying purposes.
// RecordWin (Updated) simulates recording a win.
// If lib.User doesn't exist, creates them with score 1. Otherwise increments.
func (s *SpyPlayerStore) RecordWin(name string) {
	s.recordWinCalls = append(s.recordWinCalls, name)
	// Get current score (or 0 if new) - ignore error for simplicity in mock logic here
	currentScore, _ := s.GetPlayerScore(name)
	s.scores[name] = currentScore + 1 // Increment score
}

func (s *SpyPlayerStore) GetLeague() []lib.User {

	league := make([]lib.User, 0, len(s.scores))
	for name, score := range s.scores {
		league = append(league, lib.User{Name: name, Score: score})
	}

	// Sort by score descending
	sort.Slice(league, func(a, b int) bool {
		return league[a].Score > league[b].Score
	})

	return league
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
// ErrUserNotFound definition
// lib.User struct definition
// PlayerStore interface definition
// SpyPlayerStore struct and methods (updated GetPlayerScore, RecordWin)

// PlayerServer struct definition (as before)

// Helper function to create a PlayerServer instance for testing (Updated Handlers)
func setupTestServer(t *testing.T) (*PlayerServer, *SpyPlayerStore) {
	t.Helper()
	s := NewSpyPlayerStore(t)
	mux := http.NewServeMux()

	// --- Updated GET Handler Simulation ---
	mux.HandleFunc("GET /user/{name}/score", func(w http.ResponseWriter, r *http.Request) {
		playerName := r.PathValue("name") // Requires Go 1.22+

		score, err := s.GetPlayerScore(playerName)

		if errors.Is(err, store.ErrUserNotFound) {
			// lib.User not found scenario
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			// Other potential store errors
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			// Log the actual error server-side: log.Printf("Store error: %v", err)
			return
		}

		// lib.User found, prepare JSON response
		user := lib.User{Name: playerName, Score: score}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // Explicitly set OK status

		// Marshal lib.User struct to JSON and write to response
		if err := json.NewEncoder(w).Encode(user); err != nil {
			// Handle potential JSON encoding error
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			// Log the actual error server-side: log.Printf("JSON encode error: %v", err)
		}
	})

	// --- Updated PUT Handler Simulation (Response unchanged for now) ---
	mux.HandleFunc("PUT /user/{name}/score", func(w http.ResponseWriter, r *http.Request) {
		playerName := r.PathValue("name")  // Requires Go 1.22+
		s.RecordWin(playerName)            // Call the updated RecordWin
		w.WriteHeader(http.StatusAccepted) // Keep response as Accepted
	})
	// --- End Route Setup Simulation ---

	server := &PlayerServer{
		store:   s,
		Handler: mux,
	}
	return server, s
}

// --- Updated GET Tests ---
func TestPlayerServer_GETScore(t *testing.T) {
	server, s := setupTestServer(t)

	// Define test cases for GET (expectedBody now interface{} for JSON unmarshalling)
	tests := []struct {
		name            string
		playerName      string
		initialScore    int  // Score to stub in the store
		expectUserFound bool // Flag to distinguish 404 from 200
		expectedStatus  int
		expectedBody    lib.User // Expected lib.User struct for successful cases
	}{
		{
			name:            "Get score for existing player",
			playerName:      "Alice",
			initialScore:    5,
			expectUserFound: true,
			expectedStatus:  http.StatusOK,
			expectedBody:    lib.User{Name: "Alice", Score: 5},
		},
		{
			name:            "Get score for another existing player",
			playerName:      "Bob",
			initialScore:    10,
			expectUserFound: true,
			expectedStatus:  http.StatusOK,
			expectedBody:    lib.User{Name: "Bob", Score: 10},
		},
		{
			name:            "Get score for existing player with zero score",
			playerName:      "ZeroZorro",
			initialScore:    0, // Explicitly test score 0 vs not found
			expectUserFound: true,
			expectedStatus:  http.StatusOK,
			expectedBody:    lib.User{Name: "ZeroZorro", Score: 0},
		},
		{
			name:            "Get score for non-existent player",
			playerName:      "Charlie",
			initialScore:    0,     // Store starts empty
			expectUserFound: false, // Expect a 404
			expectedStatus:  http.StatusNotFound,
			expectedBody:    lib.User{}, // Body will be empty/ignored for 404
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Stub the score only if the lib.User is expected to exist
			if tt.expectUserFound {
				s.StubScore(tt.playerName, tt.initialScore)
			}

			requestPath := fmt.Sprintf("/user/%s/score", tt.playerName)
			request, _ := http.NewRequest(http.MethodGet, requestPath, nil)
			response := httptest.NewRecorder()

			server.Handler.ServeHTTP(response, request)

			// Assert status code
			if response.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", response.Code, tt.expectedStatus)
			}

			// Assert body and Content-Type only for successful responses (200 OK)
			if tt.expectedStatus == http.StatusOK {
				// Check Content-Type header
				contentType := response.Header().Get("Content-Type")
				expectedContentType := "application/json"
				if contentType != expectedContentType {
					t.Errorf("handler returned wrong Content-Type: got %q want %q", contentType, expectedContentType)
				}

				// Decode JSON response body into a lib.User struct
				var gotUser lib.User
				err := json.NewDecoder(response.Body).Decode(&gotUser)
				if err != nil {
					t.Fatalf("Could not decode JSON response body: %v", err)
				}

				// Compare the decoded struct with the expected struct
				if !reflect.DeepEqual(gotUser, tt.expectedBody) {
					t.Errorf("handler returned unexpected body: got %+v want %+v", gotUser, tt.expectedBody)
				}
			} else if response.Body.Len() > 0 && tt.expectedStatus == http.StatusNotFound {
				// Optionally check that the body is empty on 404
				t.Errorf("expected empty body for 404, but got %q", response.Body.String())
			}
		})
	}
}

// --- Updated PUT Tests ---
func TestPlayerServer_PUTScore(t *testing.T) {
	server, s := setupTestServer(t)

	playerName := "Alice"
	requestPath := fmt.Sprintf("/user/%s/score", playerName)

	// --- First PUT (User initially doesn't exist) ---
	t.Run("first PUT creates lib.User with score 1", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPut, requestPath, nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		// Assert status code
		expectedStatus := http.StatusAccepted
		if response.Code != expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v", response.Code, expectedStatus)
		}

		// Assert store interaction
		s.AssertRecordWinCalledWith(playerName)

		// Verify score in store (using the updated GetPlayerScore)
		actualScore, err := s.GetPlayerScore(playerName)
		if err != nil {
			t.Fatalf("Error getting score after first PUT: %v", err) // Should not be ErrUserNotFound anymore
		}
		expectedScore := 1
		if actualScore != expectedScore {
			t.Errorf("score after first PUT was wrong: got %d want %d", actualScore, expectedScore)
		}
	})

	// --- Second PUT (User now exists) ---
	t.Run("second PUT increments existing lib.User score", func(t *testing.T) {
		// Ensure previous state is reflected (Alice has score 1)
		request, _ := http.NewRequest(http.MethodPut, requestPath, nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		// Assert status code
		expectedStatus := http.StatusAccepted
		if response.Code != expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v", response.Code, expectedStatus)
		}

		// Verify score in store
		actualScore, err := s.GetPlayerScore(playerName)
		if err != nil {
			t.Fatalf("Error getting score after second PUT: %v", err)
		}
		expectedScore := 2 // Score should now be 2
		if actualScore != expectedScore {
			t.Errorf("score after second PUT was wrong: got %d want %d", actualScore, expectedScore)
		}

		// Check total RecordWin calls if needed (should be 2 now across both subtests)
		// Depending on how t.Run isolates state or if the store is reset, adjust this check.
		// Assuming store state persists across t.Run in this setup:
		if len(s.recordWinCalls) != 2 {
			t.Errorf("expected RecordWin to be called twice in total, got %d calls", len(s.recordWinCalls))
		}
	})
}
