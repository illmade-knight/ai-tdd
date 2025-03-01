package gemini

import (
	"math"
	"testing"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	Radius float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

type Triangle struct {
	A float64
	B float64
	C float64
}

func TestShapes(t *testing.T) {
	testCases := struct {
		name              string
		shape             Shape
		expectedArea      float64
		expectedPerimeter float64
	}{
		{
			name:              "Circle - Radius 5",
			shape:             Circle{Radius: 5},
			expectedArea:      78.539816,
			expectedPerimeter: 31.415927,
		},
		{
			name:              "Circle - Radius 0",
			shape:             Circle{Radius: 0},
			expectedArea:      0,
			expectedPerimeter: 0,
		},
		{
			name:              "Circle - Radius 2.5",
			shape:             Circle{Radius: 2.5},
			expectedArea:      19.634954,
			expectedPerimeter: 15.707963,
		},
		{
			name:              "Rectangle - 4x6",
			shape:             Rectangle{Width: 4, Height: 6},
			expectedArea:      24,
			expectedPerimeter: 20,
		},
		{
			name:              "Rectangle - 0x6",
			shape:             Rectangle{Width: 0, Height: 6},
			expectedArea:      0,
			expectedPerimeter: 12,
		},
		{
			name:              "Rectangle - 4x0",
			shape:             Rectangle{Width: 4, Height: 0},
			expectedArea:      0,
			expectedPerimeter: 8,
		},
		{
			name:              "Triangle - Scalene",
			shape:             Triangle{A: 5, B: 6, C: 7},
			expectedArea:      14.696938,
			expectedPerimeter: 18,
		},
		{
			name:              "Triangle - Isosceles",
			shape:             Triangle{A: 5, B: 5, C: 6},
			expectedArea:      12,
			expectedPerimeter: 16,
		},
		{
			name:              "Triangle - Right-angled",
			shape:             Triangle{A: 3, B: 4, C: 5},
			expectedArea:      6,
			expectedPerimeter: 12,
		},
		{
			name:              "Triangle - Colinear",
			shape:             Triangle{A: 1, B: 2, C: 3},
			expectedArea:      0,
			expectedPerimeter: 6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			area := tc.shape.Area()
			perimeter := tc.shape.Perimeter()

			if !almostEqual(area, tc.expectedArea) {
				t.Errorf("Expected area: %f, but got: %f", tc.expectedArea, area)
			}

			if !almostEqual(perimeter, tc.expectedPerimeter) {
				t.Errorf("Expected perimeter: %f, but got: %f", tc.expectedPerimeter, perimeter)
			}
		})
	}
}

func almostEqual(a, b float64) bool {
	const tolerance = 1e-6
	return math.Abs(a-b) < tolerance
}
