Okay, let's write an integration test for the PlayerServer using the InMemoryPlayerStore. This test will interact with the server via HTTP requests, verifying the combined behavior of the server logic and the in-memory store implementation. We'll structure it to facilitate testing different store implementations later.

**Assumptions:**

1. Your PlayerServer struct exists and has a way to be instantiated with a PlayerStore.
2. The PlayerServer struct exposes its configured http.Handler (e.g., via a field named Handler) so it can be passed to httptest.NewServer. If PlayerServer has a ServeHTTP method itself, that can be used directly. We'll assume it has a Handler field based on previous test setups.
3. The InMemoryPlayerStore code from the artifact go\_inmemory\_store is available in the same package or an importable one.
4. The User struct and ErrUserNotFound are also accessible.

PlayerServer Integration Test (In-Memory)Apr 16, 11:43 AM

````go
package main_test // Or your service package_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	// Import the package containing your server and store implementation if needed
	// main "path/to/your/server/code"
)

// --- Duplicated definitions for clarity if not imported ---
// (Ideally, import these from their actual location)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type PlayerStore interface {
	GetPlayerScore(name string) (int, error)
	RecordWin(name string)
}

// Assume InMemoryPlayerStore struct and methods are defined as in go_inmemory_store artifact
type InMemoryPlayerStore struct {
	mu     sync.RWMutex
	scores map[string]int
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{scores: make(map[string]int)}
}
func (i *InMemoryPlayerStore) GetPlayerScore(name string) (int, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	score, ok := i.scores[name]
	if !ok { return 0, ErrUserNotFound }
	return score, nil
}
func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.scores[name]++
}


// Assume PlayerServer struct exists
type PlayerServer struct {
	Store   PlayerStore
	Handler http.Handler // Assumes the configured handler is accessible
}

// Assume a function exists to create and configure the server's handler
// This simulates calling server.Start() without blocking, just setting up routes.
func NewTestPlayerServer(store PlayerStore) *PlayerServer {
    mux := http.NewServeMux() // Or your router of choice

    // GET Handler logic (copied from previous test setup for completeness)
    mux.HandleFunc("GET /user/{name}/score", func(w http.ResponseWriter, r *http.Request) {
        playerName := r.PathValue("name") // Go 1.22+
        score, err := store.GetPlayerScore(playerName)
        if errors.Is(err, ErrUserNotFound) {
            w.WriteHeader(http.StatusNotFound)
            return
        } else if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        user := User{Name: playerName, Score: score}
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(user)
    })

    // PUT Handler logic (copied from previous test setup for completeness)
    mux.HandleFunc("PUT /user/{name}/score", func(w http.ResponseWriter, r *http.Request) {
        playerName := r.PathValue("name") // Go 1.22+
        store.RecordWin(playerName)
        w.WriteHeader(http.StatusAccepted)
    })

    return &PlayerServer{
        Store:   store,
        Handler: mux, // The configured router/mux is the handler
    }
}


// --- Integration Test ---

func TestPlayerServerIntegration(t *testing.T) {
	// --- Test Setup ---
	// 1. Create the store instance we want to test against.
	//    (This is where you could swap in a different store implementation later)
	store := NewInMemoryPlayerStore()

	// 2. Create the PlayerServer using the chosen store.
	//    This assumes NewTestPlayerServer correctly sets up the routes/handler.
	playerServer := NewTestPlayerServer(store)

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
		putUrl := fmt.Sprintf("%s/user/%s/score", testServer.URL, playerName)
		getUrl := fmt.Sprintf("%s/user/%s/score", testServer.URL, playerName)

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
		expectedUser := User{Name: playerName, Score: 1}
		assertUserResponse(t, getResp.Body, expectedUser)
	})

	t.Run("Record multiple wins", func(t *testing.T) {
		playerName := "PlayerTwo"
		putUrl := fmt.Sprintf("%s/user/%s/score", testServer.URL, playerName)
		getUrl := fmt.Sprintf("%s/user/%s/score", testServer.URL, playerName)

		// Record win 1
		req1, _ := http.NewRequest(http.MethodPut, putUrl, nil)
		resp1, err := testServer.Client().Do(req1)
		if err != nil { t.Fatalf("PUT 1 failed: %v", err) }
		defer resp1.Body.Close()
		assertStatusCode(t, resp1.StatusCode, http.StatusAccepted)

		// Record win 2
		req2, _ := http.NewRequest(http.MethodPut, putUrl, nil)
		resp2, err := testServer.Client().Do(req2)
		if err != nil { t.Fatalf("PUT 2 failed: %v", err) }
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
		expectedUser := User{Name: playerName, Score: 2} // Expect score 2
		assertUserResponse(t, getResp.Body, expectedUser)
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

func assertUserResponse(t *testing.T, body io.Reader, want User) {
	t.Helper()
	var got User
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

````

**Explanation:**

1. **Test Function:** TestPlayerServerIntegration serves as the main entry point.
2. **Setup:**
    * NewInMemoryPlayerStore() creates the specific store implementation we want to test against. **(Swap Point)** You can replace this line to create a different PlayerStore (e.g., connected to a test database) later.
    * NewTestPlayerServer(store) creates the server instance configured with the store. (This helper simulates the route setup part of your actual server initialization).
    * httptest.NewServer(playerServer.Handler) starts a real HTTP server on a random available port, using the handler configured within playerServer.
    * t.Cleanup(testServer.Close) ensures the test server is shut down reliably after the test function completes.
3. **Scenarios (t.Run)**: Each t.Run block tests a specific interaction flow.
4. **HTTP Client**: We use testServer.Client() which is pre-configured to work with the test server (handling cookies, etc., though not strictly needed here) or standard http.Get/http.NewRequest. URLs are constructed using testServer.URL.
5. **Assertions**: Helper functions (assertStatusCode, assertContentType, assertUserResponse) are used to check the HTTP status code, Content-Type header, and decode/compare the JSON response body against the expected User struct.
6. **State**: Because store is created once for the whole TestPlayerServerIntegration function, state persists between the t.Run blocks (e.g., PlayerOne created in the second t.Run would still exist if accessed in the third). If you need isolated state for each scenario, you would move the Setup steps *inside* each t.Run block.

This integration test provides higher confidence that your server and chosen store implementation work correctly together when communicating over HTTP.