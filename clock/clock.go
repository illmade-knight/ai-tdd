package clock

import (
	"math"
	"time"
)

type Hand int64

var (
	start       = math.Pi / 2
	circle      = math.Pi * 2
	minute Hand = 60
	hour   Hand = 12
	// we have 60 seconds and 60 minutes so the same
	second = minute
)

type Point struct {
	X float64
	Y float64
}

func hourHandPoint(t time.Time) Point {
	adjust := float64(t.Minute()) / 60.0
	hourHand := float64(t.Hour() % 12)
	return unitPoint(hourHand+adjust, hour)
}

func unitPoint(t float64, hand Hand) Point {
	return Point{math.Cos(start - t*circle/float64(hand)), math.Sin(start - t*circle/float64(hand))}
}
