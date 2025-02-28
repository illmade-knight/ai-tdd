package shapes

import (
	"testing"
)

//type Shape interface {
//	Area() float64
//	Perimeter() float64
//}
//
//type Circle struct {
//	Radius float64
//}
//
//type Rectangle struct {
//	Width  float64
//	Height float64
//}
//
//type Triangle struct {
//	A float64
//	B float64
//	C float64
//}

func V1TestShapes(t *testing.T) {
	testCases := []struct {
		name              string
		shape             Shape
		expectedArea      float64
		expectedPerimeter float64
	}{
		{
			name:              "Circle",
			shape:             Circle{Radius: 5},
			expectedArea:      78.539816, // Example value, replace with actual calculation
			expectedPerimeter: 31.415927, // Example value, replace with actual calculation
		},
		{
			name:              "Rectangle",
			shape:             Rectangle{Width: 4, Height: 6},
			expectedArea:      24, // Example value, replace with actual calculation
			expectedPerimeter: 20, // Example value, replace with actual calculation
		},
		{
			name:              "Triangle",
			shape:             Triangle{A: 3, B: 4, C: 5},
			expectedArea:      6,  // Example value, replace with actual calculation
			expectedPerimeter: 12, // Example value, replace with actual calculation
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			area := tc.shape.Area()
			perimeter := tc.shape.Perimeter()

			// Use almostEqual for floating-point comparisons if necessary
			if area != tc.expectedArea {
				t.Errorf("Expected area: %f, but got: %f", tc.expectedArea, area)
			}

			if perimeter != tc.expectedPerimeter {
				t.Errorf("Expected perimeter: %f, but got: %f", tc.expectedPerimeter, perimeter)
			}
		})
	}
}
