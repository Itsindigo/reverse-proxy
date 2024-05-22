package math_utils

import (
	"golang.org/x/exp/constraints"
	"math"
)

func Round(i float64) int {
	return int(i + math.Copysign(0.5, i))
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
