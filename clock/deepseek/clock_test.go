package deepseek

import (
	"clock"
	"math"
	"testing"
	"time"
)

// Replace GetHands with the actual function from your library.
// func GetHands(t time.Time) (hour, minute, second clock.Point) {
//     // Implementation
// }

func TestHandPositions(t *testing.T) {
	tests := []struct {
		name    string
		timeStr string // Format "2006-01-02 15:04:05"
		hour    clock.Point
		minute  clock.Point
		second  clock.Point
	}{
		{
			name:    "12:00:00",
			timeStr: "2023-10-23 12:00:00",
			hour:    clock.Point{X: 0, Y: -1},
			minute:  clock.Point{X: 0, Y: -1},
			second:  clock.Point{X: 0, Y: -1},
		},
		{
			name:    "3:00:00",
			timeStr: "2023-10-23 03:00:00",
			hour:    clock.Point{X: 1, Y: 0},
			minute:  clock.Point{X: 0, Y: -1},
			second:  clock.Point{X: 0, Y: -1},
		},
		{
			name:    "6:00:00",
			timeStr: "2023-10-23 06:00:00",
			hour:    clock.Point{X: 0, Y: 1},
			minute:  clock.Point{X: 0, Y: -1},
			second:  clock.Point{X: 0, Y: -1},
		},
		{
			name:    "9:30:15",
			timeStr: "2023-10-23 09:30:15",
			hour:    clock.Point{X: -0.9659258, Y: -0.2588190},
			minute:  clock.Point{X: -0.0261799, Y: 0.9996573},
			second:  clock.Point{X: 1, Y: 0},
		},
		{
			name:    "12:30:45",
			timeStr: "2023-10-23 12:30:45",
			hour:    clock.Point{X: -0.1305262, Y: -0.9914448},
			minute:  clock.Point{X: -0.7071068, Y: -0.7071068},
			second:  clock.Point{X: -0.8660254, Y: -0.5},
		},
	}

	epsilon := 1e-6
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm, err := time.Parse("2006-01-02 15:04:05", tt.timeStr)
			if err != nil {
				t.Fatalf("could not parse time: %v", err)
			}
			hands := GetHands(tm)
			hour, minute, second := hands.Hour, hands.Minute, hands.Second

			if !approxEqual(hour, tt.hour, epsilon) {
				t.Errorf("hour hand position incorrect: expected %+v, got %+v", tt.hour, hour)
			}
			if !approxEqual(minute, tt.minute, epsilon) {
				t.Errorf("minute hand position incorrect: expected %+v, got %+v", tt.minute, minute)
			}
			if !approxEqual(second, tt.second, epsilon) {
				t.Errorf("second hand position incorrect: expected %+v, got %+v", tt.second, second)
			}
		})
	}
}

func approxEqual(a, b clock.Point, epsilon float64) bool {
	return math.Abs(a.X-b.X) <= epsilon && math.Abs(a.Y-b.Y) <= epsilon
}
