package basic

import "strings"

// Function to be tested
func RepeatCharacter(char rune, count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat(string(char), count)
}
