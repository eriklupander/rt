package mat

import "math"

const Epsilon = 0.001

func Eq(a, b float64) bool {
	return math.Abs(a-b) < Epsilon
}
