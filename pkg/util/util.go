package util

import "math"

func RoundFloat64(f float64) float64 {
	return math.Floor(f + .5)
}
