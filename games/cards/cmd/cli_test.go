package cmd

import (
	"games/user/store"
	"strings"
	"testing"
)

func assertPlayerWins(t testing.TB, s *store.InMemoryPlayerStore, winner string) {
	t.Helper()

	winnerScore, err := s.GetPlayerScore(winner)
	if err != nil {
		t.Fatal(err)
	}
	if winnerScore == 0 {
		t.Fatalf("got %d wins want %d", winnerScore, 1)
	}

}

func TestCLI(t *testing.T) {
	playerStore := store.NewInMemoryPlayerStore()
	in := strings.NewReader("Chris wins\n")
	cli := &CLI{playerStore, in}
	cli.PlayCards()

	assertPlayerWins(t, playerStore, "Chris")
}
