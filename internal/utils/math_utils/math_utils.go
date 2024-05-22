package math_utils

import (
	"math"
)

func Round(i float64) int {
	return int(i + math.Copysign(0.5, i))
}

func ToFixed(i float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(i*output)) / output
}
