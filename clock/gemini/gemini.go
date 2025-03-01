package gemini

import (
	"fmt"
	"math"
	"testing"
	"time"
)

// A Point represents a two dimensional Cartesian coordinate
type Point struct {
	X float64
	Y float64
}

type Hand int

const (
	second Hand = 60
	minute Hand = 60
	hour   Hand = 12
)

func TestSecondHandAtMidnight(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 0, 0, 0, 0, time.UTC)

	want := Point{X: 0, Y: 1} // 12 o'clock
	got := SecondHand(tm)

	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}

func TestSecondHandAt30Seconds(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 0, 0, 30, 0, time.UTC)

	want := Point{X: 0, Y: -1} // 6 o'clock
	got := SecondHand(tm)

	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}

// ... similar tests for other SecondHand positions

func TestMinuteHandAtMidnight(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 0, 0, 0, 0, time.UTC)

	want := Point{X: 0, Y: 1} // 12 o'clock
	got := MinuteHand(tm)

	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}

func TestMinuteHandAt15Minutes(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 0, 15, 0, 0, time.UTC)

	want := Point{X: -1, Y: 0} // 3 o'clock
	got := MinuteHand(tm)

	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}

// ... similar tests for other MinuteHand positions

func TestHourHandAtMidnight(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 0, 0, 0, 0, time.UTC)

	want := Point{X: 0, Y: 1} // 12 o'clock
	got := HourHand(tm)

	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}

func TestHourHandAt3AM(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 3, 0, 0, 0, time.UTC)

	want := Point{X: -1, Y: 0} // 3 o'clock
	got := HourHand(tm)

	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}

// ... similar tests for other HourHand positions

func TestUnitPoint(t *testing.T) {
	testCases := struct {
		hand  Hand
		value float64
		want  Point
	}{
		{second, 0, Point{0, 1}},
		{second, 15, Point{-1, 0}},
		{second, 30, Point{0, -1}},
		{second, 45, Point{1, 0}},

		{minute, 0, Point{0, 1}},
		{minute, 15, Point{-1, 0}},
		{minute, 30, Point{0, -1}},
		{minute, 45, Point{1, 0}},

		{hour, 0, Point{0, 1}},
		{hour, 3, Point{-1, 0}},
		{hour, 6, Point{0, -1}},
		{hour, 9, Point{1, 0}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v at %v", tc.hand, tc.value), func(t *testing.T) {
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
