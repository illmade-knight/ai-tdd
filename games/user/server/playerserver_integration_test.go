package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"games/user/lib"
	"games/user/store"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"testing"
)

// Assume a function exists to create and configure the server's handler
// This simulates calling server.Start() without blocking, just setting up routes.
func NewTestPlayerServer(s store.PlayerStore) *PlayerServer {
	//mux := http.NewServeMux() // Or your router of choice
	//
	//// GET Handler logic (copied from previous test setup for completeness)
	//mux.HandleFunc("GET /user/{name}/score", func(w http.ResponseWriter, r *http.Request) {
	//	playerName := r.PathValue("name") // Go 1.22+
	//	score, err := s.GetPlayerScore(playerName)
	//	if errors.Is(err, store.ErrUserNotFound) {
	//		w.WriteHeader(http.StatusNotFound)
	//		return
	//	} else if err != nil {
	//		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	//		return
	//	}
	//	user := lib.User{Name: playerName, Score: score}
	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(http.StatusOK)
	//	json.NewEncoder(w).Encode(user)
	//})
	//
	//// PUT Handler logic (copied from previous test setup for completeness)
	//mux.HandleFunc("PUT /user/{name}/score", func(w http.ResponseWriter, r *http.Request) {
	//	playerName := r.PathValue("name") // Go 1.22+
	//	s.RecordWin(playerName)
	//	w.WriteHeader(http.StatusAccepted)
	//})
	//
	//mux.HandleFunc("GET /league", func(w http.ResponseWriter, r *http.Request) {
	//
	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(http.StatusOK)
	//
	//	league, err := s.GetLeague()
	//	if err != nil {
	//		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	//	}
	//	json.NewEncoder(w).Encode(league)
	//})
	//
	//return &PlayerServer{
	//	store:   s,
	//	Handler: mux, // The configured router/mux is the handler
	//}
	return NewPlayerServer(store.NewInMemoryPlayerStore())
}

// --- Integration Test ---

func TestPlayerServerIntegration(t *testing.T) {
	// --- Test Setup ---
	// 1. Create the store instance we want to test against.
	//    (This is where you could swap in a different store implementation later)
	s := store.NewInMemoryPlayerStore()

	// 2. Create the PlayerServer using the chosen store.
	//    This assumes NewTestPlayerServer correctly sets up the routes/handler.
	playerServer := NewTestPlayerServer(s)

	// 3. Start a test HTTP server using the PlayerServer's handler.
	//    httptest.NewServer finds an available port and listens on it.
	testServer := httptest.NewServer(playerServer.Handler)

	// 4. Register cleanup function to close the server when the test finishes.
	//    This is crucial to release the network port.
	t.Cleanup(func() {
		testServer.Close()
	})

	// --- Test Scenarios ---

	t.Run("Get score for non-existent player returns 404", func(t *testing.T) {
		// Construct the request URL using the test server's dynamic URL
		url := fmt.Sprintf("%s/user/NonExistentUser/score", testServer.URL)

		// Send GET request
		resp, err := http.Get(url)
		if err != nil {
			t.Fatalf("Could not send GET request: %v", err)
		}
		defer resp.Body.Close() // Ensure body is closed

		// Assert status code
		assertStatusCode(t, resp.StatusCode, http.StatusNotFound)
	})

	t.Run("Record win and get score", func(t *testing.T) {
		playerName := "PlayerOne"
		putUrl := fmt.Sprintf("%s/players/%s/wins", testServer.URL, playerName)
		getUrl := fmt.Sprintf("%s/players/%s/wins", testServer.URL, playerName)

		// 1. Record a win (PUT request)
		req, err := http.NewRequest(http.MethodPut, putUrl, nil)
		if err != nil {
			t.Fatalf("Could not create PUT request: %v", err)
		}
		// Use the test server's client to handle potential redirects etc.
		putResp, err := testServer.Client().Do(req)
		if err != nil {
			t.Fatalf("Could not send PUT request: %v", err)
		}
		defer putResp.Body.Close()

		// Assert PUT status code
		assertStatusCode(t, putResp.StatusCode, http.StatusAccepted)

		// 2. Get the score (GET request)
		getResp, err := testServer.Client().Get(getUrl) // Use client for consistency
		if err != nil {
			t.Fatalf("Could not send GET request: %v", err)
		}
		defer getResp.Body.Close()

		// Assert GET status code
		assertStatusCode(t, getResp.StatusCode, http.StatusOK)

		// Assert Content-Type
		assertContentType(t, getResp, "application/json")

		// Assert response body
		expectedUser := lib.User{Name: playerName, Score: 1}
		assertUserResponse(t, getResp.Body, expectedUser)
	})

	t.Run("Record multiple wins", func(t *testing.T) {
		playerName := "PlayerTwo"
		putUrl := fmt.Sprintf("%s/players/%s/wins", testServer.URL, playerName)
		getUrl := fmt.Sprintf("%s/players/%s/wins", testServer.URL, playerName)

		// Record win 1
		req1, _ := http.NewRequest(http.MethodPut, putUrl, nil)
		resp1, err := testServer.Client().Do(req1)
		if err != nil {
			t.Fatalf("PUT 1 failed: %v", err)
		}
		defer resp1.Body.Close()
		assertStatusCode(t, resp1.StatusCode, http.StatusAccepted)

		// Record win 2
		req2, _ := http.NewRequest(http.MethodPut, putUrl, nil)
		resp2, err := testServer.Client().Do(req2)
		if err != nil {
			t.Fatalf("PUT 2 failed: %v", err)
		}
		defer resp2.Body.Close()
		assertStatusCode(t, resp2.StatusCode, http.StatusAccepted)

		// Get final score
		getResp, err := testServer.Client().Get(getUrl)
		if err != nil {
			t.Fatalf("Could not send GET request: %v", err)
		}
		defer getResp.Body.Close()

		// Assert status and body
		assertStatusCode(t, getResp.StatusCode, http.StatusOK)
		assertContentType(t, getResp, "application/json")
		expectedUser := lib.User{Name: playerName, Score: 2} // Expect score 2
		assertUserResponse(t, getResp.Body, expectedUser)
	})

	// Helper function for setting up server and store for a subtest
	setup := func(t *testing.T) (*httptest.Server, store.PlayerStore) {
		t.Helper()
		//playerStore := store.NewInMemoryPlayerStore()

		playerStore, _ := createTempParquetStore(t)

		playerServer := NewTestPlayerServer(playerStore)
		testServer := httptest.NewServer(playerServer.Handler)
		t.Cleanup(func() { testServer.Close() })
		return testServer, playerStore
	}
	// --- New Failing Test for /league ---
	t.Run("Get league table", func(t *testing.T) {
		// Setup server and store specifically for this test
		testServer, s := setup(t)

		// Pre-populate the store with data for the league table
		s.RecordWin("Charlie") // Score 1
		s.RecordWin("Alice")   // Score 1
		s.RecordWin("Alice")   // Score 2
		s.RecordWin("Bob")     // Score 1
		s.RecordWin("Bob")     // Score 2
		s.RecordWin("Bob")     // Score 3

		// Expected league order: Bob (3), Alice (2), Charlie (1)
		expectedLeague := []lib.User{
			{Name: "Bob", Score: 3},
			{Name: "Alice", Score: 2},
			{Name: "Charlie", Score: 1},
		}

		// Construct request URL
		url := fmt.Sprintf("%s/league", testServer.URL)

		// Send GET request to the league endpoint
		resp, err := testServer.Client().Get(url)
		if err != nil {
			t.Fatalf("Could not send GET request to /league: %v", err)
		}
		defer resp.Body.Close()

		// Assertions (These will fail initially)
		assertStatusCode(t, resp.StatusCode, http.StatusOK) // Will fail (likely 404)
		assertContentType(t, resp, "application/json")      // Will fail

		// Assert response body matches the expected sorted league
		var gotLeague []lib.User
		err = json.NewDecoder(resp.Body).Decode(&gotLeague)
		if err != nil {
			// Read body for error context if decode fails
			bodyBytes, _ := io.ReadAll(io.MultiReader(bytes.NewReader([]byte{}), resp.Body))
			t.Fatalf("Could not decode JSON league response body: %v. Raw body: %s", err, string(bodyBytes))
		}

		if !reflect.DeepEqual(gotLeague, expectedLeague) { // Will fail
			t.Errorf("handler returned unexpected league body: got %+v want %+v", gotLeague, expectedLeague)
		}
	})
}

// --- Helper Assertions ---

func assertStatusCode(t *testing.T, got, want int) {
	t.Helper() // Mark as test helper
	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}

func assertContentType(t *testing.T, resp *http.Response, want string) {
	t.Helper()
	if got := resp.Header.Get("Content-Type"); got != want {
		t.Errorf("handler returned wrong Content-Type: got %q want %q", got, want)
	}
}

func assertUserResponse(t *testing.T, body io.Reader, want lib.User) {
	t.Helper()
	var got lib.User
	err := json.NewDecoder(body).Decode(&got)
	if err != nil {
		// Try reading the raw body for better error message if decode fails
		bodyBytes, _ := io.ReadAll(io.MultiReader( // Need to reconstruct reader if already read
			bytes.NewReader([]byte{}), // Placeholder if body was fully consumed
			body,                      // Original body
		))
		t.Fatalf("Could not decode JSON response body: %v. Raw body: %s", err, string(bodyBytes))
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("handler returned unexpected body: got %+v want %+v", got, want)
	}
}

// createTempParquetStore creates a ParquetPlayerStore using a temporary file.
// It returns the store instance and the path to the temp file.
// It uses t.Fatalf on errors during setup.
func createTempParquetStore(t *testing.T) (*store.ParquetPlayerStore, string) {
	t.Helper() // Mark as test helper

	tempDir := t.TempDir() // Create temp dir, automatically cleaned up
	tempFilePath := filepath.Join(tempDir, "test_players.parquet")

	// Ensure NewParquetPlayerStore is accessible (might need import or be in same package)
	s, err := store.NewParquetPlayerStore(tempFilePath)
	if err != nil {
		t.Fatalf("Failed to create ParquetPlayerStore for test: %v", err)
	}
	return s, tempFilePath
}
