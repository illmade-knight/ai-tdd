package basic

import (
	"testing"
)

func TestRepeatCharacter(t *testing.T) {
	testCases := []struct {
		name     string
		char     rune
		count    int
		expected string
	}{
		{
			name:     "Repeat 'a' 5 times",
			char:     'a',
			count:    5,
			expected: "aaaaa",
		},
		{
			name:     "Repeat 'B' 0 times",
			char:     'B',
			count:    0,
			expected: "",
		},
		{
			name:     "Repeat '$' 1 time",
			char:     '$',
			count:    1,
			expected: "$",
		},
		{
			name:     "Repeat ' ' 3 times", // Space character
			char:     ' ',
			count:    3,
			expected: "   ",
		},
		{
			name:     "Repeat unicode character '世' twice",
			char:     '世',
			count:    2,
			expected: "世世",
		},
		{
			name:     "Repeat 'z' -1 times", // Test negative case
			char:     'z',
			count:    -1,
			expected: "", // Expect empty string or your preferred error handling.
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := RepeatCharacter(tc.char, tc.count)
			if result != tc.expected {
				t.Errorf("Expected: %q, but got: %q for char: %c, count: %d", tc.expected, result, tc.char, tc.count)
			}
		})
	}
}
