package clock

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestHourHandVector(t *testing.T) {

	cases := []struct {
		time  time.Time
		point Point
	}{
		{simpleTime(6, 0, 0), Point{0, -1}},
		{simpleTime(6, 30, 0), Point{-0.2588190451025204, -0.9659258262890684}},
		{simpleTime(18, 30, 0), Point{-0.2588190451025204, -0.9659258262890684}},
		{simpleTime(21, 0, 0), Point{-1, 0}},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := hourHandPoint(c.time)
			if !roughlyEqualPoint(got, c.point) {
				t.Fatalf("Wanted %v Point, but got %v", c.point, got)
			}
		})
	}
}

func TestUnitPoint(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		hand  Hand
		value float64
		want  Point
	}{
		// Second Hand test cases
		{second, 0, Point{0, 1}},   // 12 o'clock
		{second, 15, Point{1, 0}},  // 3 o'clock
		{second, 30, Point{0, -1}}, // 6 o'clock
		{second, 45, Point{-1, 0}}, // 9 o'clock

		// Minute Hand test cases
		{minute, 0, Point{0, 1}},
		{minute, 15, Point{1, 0}},
		{minute, 30, Point{0, -1}},
		{minute, 45, Point{-1, 0}},

		// Hour Hand test cases
		{hour, 0, Point{0, 1}},
		{hour, 3, Point{1, 0}},
		{hour, 6, Point{0, -1}},
		{hour, 9, Point{-1, 0}},
	}

	// Run the tests
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Hand %v at %v", tc.hand, tc.value), func(t *testing.T) {
			got := unitPoint(tc.value, tc.hand)
			if !roughlyEqualPoint(got, tc.want) {
				t.Errorf("Got %v, wanted %v", got, tc.want)
			}
		})
	}
}

func simpleTime(hours, minutes, seconds int) time.Time {
	return time.Date(312, time.October, 28, hours, minutes, seconds, 0, time.UTC)
}

func testName(t time.Time) string {
	return t.Format("15:04:05")
}

func roughlyEqualFloat64(a, b float64) bool {
	const equalityThreshold = 1e-7
	return math.Abs(a-b) < equalityThreshold
}

func roughlyEqualPoint(a, b Point) bool {
	return roughlyEqualFloat64(a.X, b.X) &&
		roughlyEqualFloat64(a.Y, b.Y)
}
