package basic

import "testing"

func TestSum(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{
			name:     "Sum of no numbers",
			input:    []int{},
			expected: 0,
		},
		{
			name:     "Sum of one number, single int",
			input:    5,
			expected: 5,
		},
		{
			name:     "Sum of one number",
			input:    []int{5},
			expected: 5,
		},
		{
			name:     "Sum of two numbers",
			input:    []int{2, 3},
			expected: 5,
		},
		{
			name:     "Sum of multiple numbers",
			input:    []int{1, 2, 3, 4, 5},
			expected: 15,
		},
		{
			name:     "Sum of negative numbers",
			input:    []int{-1, -2, -3},
			expected: -6,
		},
		{
			name:     "Sum of mixed positive and negative numbers",
			input:    []int{-1, 2, -3, 4, -5, 6},
			expected: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Sum(toIntSlice(tc.input)...) // Important: Use ... to unpack the slice

			if result != tc.expected {
				t.Errorf("Expected: %d, but got: %d for input: %v", tc.expected, result, tc.input)
			}
		})
	}
}

func toIntSlice(input interface{}) []int {
	switch v := input.(type) {
	case []int:
		return v
	case int:
		return []int{v} // Create a slice with the single int
	default:
		panic("unsupported type") // Or handle other types if needed
	}
}

func TestSumMixed(t *testing.T) {
	testCases := []struct {
		name        string
		input       []interface{} // Now a slice of interface{}
		expectedSum int
		expectedErr string // Expected error message (empty string if no error)
	}{
		{
			name:        "Empty input",
			input:       []interface{}{},
			expectedSum: 0,
			expectedErr: "",
		},
		{
			name:        "Single int",
			input:       []interface{}{5},
			expectedSum: 5,
			expectedErr: "",
		},
		{
			name:        "Single slice",
			input:       []interface{}{[]int{1, 2, 3}},
			expectedSum: 6,
			expectedErr: "",
		},
		{
			name:        "Mixed ints and slices",
			input:       []interface{}{1, []int{2, 3}, 4, []int{5, 6}},
			expectedSum: 21,
			expectedErr: "",
		},
		{
			name:        "Mixed ints, slices, and negative numbers",
			input:       []interface{}{1, []int{-2, 3}, -4, []int{5, -6}},
			expectedSum: -3,
			expectedErr: "",
		},
		{
			name:        "Invalid input type",
			input:       []interface{}{"hello"}, // String is invalid
			expectedSum: 0,
			expectedErr: "invalid input type: string",
		},
		{
			name:        "Mixed valid and invalid types",
			input:       []interface{}{1, "hello", []int{2, 3}}, // String is invalid
			expectedSum: 6,
			expectedErr: "invalid input type: string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sum, err := SumMixed(tc.input...)
			if sum != tc.expectedSum {
				t.Errorf("Expected sum: %d, but got: %d for input: %v", tc.expectedSum, sum, tc.input)
			}

			if tc.expectedErr == "" {
				if err != nil {
					t.Errorf("Expected no error, but got: %v for input: %v", err, tc.input)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error: %q, but got nil for input: %v", tc.expectedErr, tc.input)
				} else if err.Error() != tc.expectedErr {
					t.Errorf("Expected error: %q, but got: %q for input: %v", tc.expectedErr, err.Error(), tc.input)
				}
			}
		})
	}
}
