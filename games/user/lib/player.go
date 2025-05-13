package lib

// User represents the data structure for a player.
type User struct {
	// Use json struct tags for correct marshalling
	Name  string `json:"name"`
	Score int    `json:"score"`
}
