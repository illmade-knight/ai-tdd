package deepseek

import (
	"clock"
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

type Hands struct {
	Hour   clock.Point
	Minute clock.Point
	Second clock.Point
}

func GetHands(t time.Time) Hands {
	h := hourHandPoint(t)
	m := unitPoint(float64(t.Minute()), minute)
	s := unitPoint(float64(t.Second()), minute)
	return Hands{
		Hour:   h,
		Minute: m,
		Second: s,
	}
}

func hourHandPoint(t time.Time) clock.Point {
	adjust := float64(t.Minute()) / 60.0
	hourHand := float64(t.Hour() % 12)
	return unitPoint(hourHand+adjust, hour)
}

func unitPoint(t float64, hand Hand) clock.Point {
	return clock.Point{math.Cos(start - t*circle/float64(hand)), math.Sin(start - t*circle/float64(hand))}
}
