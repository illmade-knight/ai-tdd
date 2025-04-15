package ai_core

import (
	"fmt"
	"testing"
)

// Define the sum function for integers, matching the expected signature for Reduce.
// It takes the accumulator (current sum) and the current element, returning the new sum.
func sumInts(accumulator int, currentItem int) int {
	// Basic sum. We can consider potential overflow later if required.
	return accumulator + currentItem
}

// TestReduceIntSum tests the Reduce function specifically for summing integers.
func TestReduceIntSum(t *testing.T) {

	// Define test cases using table-driven tests, a common Go pattern.
	testCases := []struct {
		name         string             // Descriptive name for the test case
		inputSlice   []int              // The slice to reduce
		initialValue int                // The starting value for the reduction
		sumFunc      func(int, int) int // The function to apply
		expected     int                // The expected result
	}{
		{
			name:         "Sum of positive integers",
			inputSlice:   []int{1, 2, 3, 4, 5},
			initialValue: 0, // Start summing from 0
			sumFunc:      sumInts,
			expected:     15, // 1+2+3+4+5 = 15
		},
		{
			name:         "Sum with a non-zero initial value",
			inputSlice:   []int{1, 2, 3},
			initialValue: 10, // Start summing from 10
			sumFunc:      sumInts,
			expected:     16, // 10 + 1 + 2 + 3 = 16
		},
		{
			name:         "Sum including negative numbers",
			inputSlice:   []int{10, -2, 5, -8, 3},
			initialValue: 0,
			sumFunc:      sumInts,
			expected:     8, // 10 - 2 + 5 - 8 + 3 = 8
		},
		{
			name:         "Sum of a single element slice",
			inputSlice:   []int{42},
			initialValue: 0,
			sumFunc:      sumInts,
			expected:     42, // 0 + 42 = 42
		},
		{
			name:         "Sum of an empty slice",
			inputSlice:   []int{}, // Empty slice
			initialValue: 0,
			sumFunc:      sumInts,
			expected:     0, // Reducing an empty slice should return the initial value
		},
		{
			name:         "Sum of an empty slice with non-zero initial value",
			inputSlice:   []int{},
			initialValue: 100,
			sumFunc:      sumInts,
			expected:     100, // Should still return the initial value
		},
		{
			name:         "Sum with a nil slice",
			inputSlice:   nil, // Test nil explicitly
			initialValue: 0,
			sumFunc:      sumInts,
			expected:     0, // Assume Reduce handles nil slice gracefully like an empty slice
		},
		// Potential Future Edge Case (Illustrative - Depends on requirements)
		// {
		//  name:         "Sum causing potential overflow (requires specific handling)",
		//  inputSlice:   []int{math.MaxInt64, 1},
		//  initialValue: 0,
		//  sumFunc:      sumInts,
		//  // Expected result depends on how overflow is handled (panic, wrap, error?)
		//  // expected:     ?,
		// },
	}

	// Iterate through the test cases
	for _, tc := range testCases {
		// t.Run allows running sub-tests, making output clearer
		t.Run(tc.name, func(t *testing.T) {
			// --- THIS IS WHERE YOU CALL THE ACTUAL Reduce FUNCTION ---
			// Placeholder: Replace 'functional.Reduce' with the actual call once implemented.
			// got := functional.Reduce(tc.inputSlice, tc.initialValue, tc.sumFunc)

			// For now, to make the test file compile, we simulate the call
			// and result conceptually. Remove this placeholder section later.
			// ---- START PLACEHOLDER ----
			fmt.Printf("Simulating call: Reduce(%v, %d, sumFunc) for test '%s'\n", tc.inputSlice, tc.initialValue, tc.name)
			// Simulate the logic simply for test structure validation
			placeholderResult := tc.initialValue
			if tc.inputSlice != nil { // Basic nil check for simulation
				for _, item := range tc.inputSlice {
					placeholderResult = tc.sumFunc(placeholderResult, item)
				}
			}
			got := placeholderResult
			// ---- END PLACEHOLDER ----

			// Assert the result
			if got != tc.expected {
				t.Errorf("Reduce(%v, %d, sumFunc) = %d; want %d", tc.inputSlice, tc.initialValue, got, tc.expected)
			}
		})
	}
}

// --- Further Tests to Consider for Reduce ---
// - Test with different data types (e.g., float64, string concatenation) once Reduce is generic.
// - Test ReduceRight (if implemented).
// - Test scenarios where the reduction function itself might have side effects or errors (how should Reduce handle this?).
// - Test for concurrency safety if the Reduce implementation is intended to be used concurrently.
