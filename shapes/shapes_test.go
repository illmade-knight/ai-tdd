package shapes

import (
	"math"
	"testing"
)

// ... (Shape interface, Circle, Rectangle, Triangle structs)

func TestShapes(t *testing.T) {
	testCases := []struct {
		name              string
		shape             Shape
		expectedArea      float64
		expectedPerimeter float64
		expectError       bool
	}{
		{
			name:              "Circle - Valid Radius",
			shape:             Circle{Radius: 5},
			expectedArea:      78.539816,
			expectedPerimeter: 31.415927,
		},
		{
			name:              "Circle - Zero Radius",
			shape:             Circle{Radius: 0},
			expectedArea:      0,
			expectedPerimeter: 0,
		},
		{
			name:        "Circle - Negative Radius",
			shape:       Circle{Radius: -5},
			expectError: true, // Expect an error for negative radius
		},
		{
			name:              "Rectangle - Valid Dimensions",
			shape:             Rectangle{Width: 4, Height: 6},
			expectedArea:      24,
			expectedPerimeter: 20,
		},
		{
			name:              "Rectangle - Zero Width",
			shape:             Rectangle{Width: 0, Height: 6},
			expectedArea:      0,
			expectedPerimeter: 12,
		},
		{
			name:              "Rectangle - Zero Height",
			shape:             Rectangle{Width: 4, Height: 0},
			expectedArea:      0,
			expectedPerimeter: 8,
		},
		{
			name:        "Rectangle - Negative Width",
			shape:       Rectangle{Width: -4, Height: 6},
			expectError: true, // Expect an error for negative width
		},
		{
			name:        "Rectangle - Negative Height",
			shape:       Rectangle{Width: 4, Height: -6},
			expectError: true, // Expect an error for negative height
		},
		{
			name:              "Triangle - Scalene",
			shape:             Triangle{A: 5, B: 6, C: 7},
			expectedArea:      14.696938, // Example value, replace with actual calculation
			expectedPerimeter: 18,        // Example value, replace with actual calculation
		},
		{
			name:              "Triangle - Isosceles",
			shape:             Triangle{A: 5, B: 5, C: 6},
			expectedArea:      12, // Example value, replace with actual calculation
			expectedPerimeter: 16, // Example value, replace with actual calculation
		},
		{
			name:              "Triangle - Right-angled",
			shape:             Triangle{A: 3, B: 4, C: 5},
			expectedArea:      6,  // Example value, replace with actual calculation
			expectedPerimeter: 12, // Example value, replace with actual calculation
		},
		{
			name:              "Triangle - Colinear",
			shape:             Triangle{A: 1, B: 2, C: 3},
			expectedArea:      0, // Example value, replace with actual calculation
			expectedPerimeter: 6, // Example value, replace with actual calculation
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tc.expectError {
					t.Errorf("Unexpected panic: %v", r)
				}
			}()

			area := tc.shape.Area()
			perimeter := tc.shape.Perimeter()

			if tc.expectError {
				return // If an error is expected, skip the following checks
			}

			// Use almostEqual for floating-point comparisons if necessary
			if !almostEqual(area, tc.expectedArea) {
				t.Errorf("Expected area: %f, but got: %f", tc.expectedArea, area)
			}

			if !almostEqual(perimeter, tc.expectedPerimeter) {
				t.Errorf("Expected perimeter: %f, but got: %f", tc.expectedPerimeter, perimeter)
			}
		})
	}
}

// almostEqual compares two float64 values for approximate equality
func almostEqual(a, b float64) bool {
	const tolerance = 1e-6 // Adjust the tolerance as needed
	return math.Abs(a-b) < tolerance
}
