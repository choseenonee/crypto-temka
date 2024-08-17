package utils

import "math"

func FloatEquals(a, b float64) bool {
	epsilon := 1e-9
	return math.Abs(a-b) < epsilon
}
