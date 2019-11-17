package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShearing(t *testing.T) {
	tc := []struct {
		shear []float64
		point Tuple4
		out   Tuple4
	}{
		{[]float64{1, 0, 0, 0, 0, 0}, *NewPoint(2, 3, 4), *NewPoint(5, 3, 4)},
		{[]float64{0, 1, 0, 0, 0, 0}, *NewPoint(2, 3, 4), *NewPoint(6, 3, 4)},
		{[]float64{0, 0, 1, 0, 0, 0}, *NewPoint(2, 3, 4), *NewPoint(2, 5, 4)},
		{[]float64{0, 0, 0, 1, 0, 0}, *NewPoint(2, 3, 4), *NewPoint(2, 7, 4)},
		{[]float64{0, 0, 0, 0, 1, 0}, *NewPoint(2, 3, 4), *NewPoint(2, 3, 6)},
		{[]float64{0, 0, 0, 0, 0, 1}, *NewPoint(2, 3, 4), *NewPoint(2, 3, 7)},
	}

	for _, test := range tc {
		p2 := MultiplyByTuple(*ShearBy(test.shear), test.point)
		assert.True(t, TupleEquals(test.out, *p2))
	}
}
