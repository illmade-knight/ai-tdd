package shapes

import "math"

// Define the shape interfaces and structs

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Triangle struct {
	A float64
	B float64
	C float64
}

func (t Triangle) Area() float64 {
	s := t.Perimeter() / 2                                  // Semi-perimeter
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C)) // Heron's formula
}

func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}
