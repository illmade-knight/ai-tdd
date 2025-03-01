package gemini

import (
	"sync"
	"testing"
)

type GeminiCounter interface {
	Inc()
	Value() int
}

func TestCounter(t *testing.T) {
	t.Run("Concurrent increments", func(t *testing.T) {
		counter := NewCounter() // Replace with your Counter constructor

		numRoutines := 1000
		numIncrements := 1000

		var wg sync.WaitGroup
		wg.Add(numRoutines)

		for i := 0; i < numRoutines; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < numIncrements; j++ {
					counter.Inc()
				}
			}()
		}

		wg.Wait()

		expectedCount := numRoutines * numIncrements
		if counter.Value() != expectedCount {
			t.Errorf("Expected count: %d, but got: %d", expectedCount, counter.Value())
		}
	})
}
