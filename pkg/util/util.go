package util

import "math"

func RoundFloat64(f float64) float64 {
	return math.Floor(f + .5)
}

func MaxUint8(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}
