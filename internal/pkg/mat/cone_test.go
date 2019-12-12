package mat

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestConeRayHits(t *testing.T) {
	c := NewCone()

	tc := []conetest{
		//{p: NewPoint(0, 0, -5), v: Normalize(NewVector(0, 0, 1)), t1: 5, t2: 5},
		{p: NewPoint(0, 0, -5), v: Normalize(NewVector(1, 1, 1)), t1: 8.660254037844386, t2: 8.660254037844386},
		{p: NewPoint(1, 1, -5), v: Normalize(NewVector(-0.5, -1, 1)), t1: 4.550055679356349, t2: 49.449944320643645},
	}
	for _, test := range tc {
		r := NewRay(test.p, test.v)
		xs := c.IntersectLocal(r)
		assert.Len(t, xs, 2)
		assert.Equal(t, test.t1, xs[0].T)
		assert.Equal(t, test.t2, xs[1].T)
	}
}

func TestIntersectConeWithParallellRay(t *testing.T) {
	c := NewCone()
	direction := Normalize(NewVector(0, 1, 1))
	r := NewRay(NewPoint(0, 0, -1), direction)
	xs := c.IntersectLocal(r)
	assert.Equal(t, 1, len(xs))
	assert.Equal(t, 0.3535533905932738, xs[0].T)
}

func TestIntersectConeCap(t *testing.T) {
	cone := NewConeMMC(-0.5, 0.5, true)

	tc := []conetest{
		//{p: NewPoint(0, 0, -5), v: Normalize(NewVector(0, 1, 0)), t1: 0,},
		//{p: NewPoint(0, 0, -0.25), v:  Normalize(NewVector(0, 1, 1)), t1:  2},
		{p: NewPoint(0, 0, -0.25), v: Normalize(NewVector(0, 1, 0)), t1: 4},
	}
	for _, test := range tc {
		r := NewRay(test.p, test.v)
		xs := cone.IntersectLocal(r)
		assert.Equal(t, int(test.t1), len(xs))
	}
}

func TestNormalOnCone(t *testing.T) {
	cone := NewConeMMC(-0.5, 0.5, true)

	tc := []conetest{
		{p: NewPoint(0, 0, 0), v: NewVector(0, 0, 0)},
		{p: NewPoint(1, 1, 1), v: NewVector(1, -math.Sqrt(2), 1)},
		{p: NewPoint(-1, -1, 0), v: NewVector(-1, 1, 0)},
	}
	for _, test := range tc {
		n := cone.NormalAtLocal(test.p, nil)
		assert.Equal(t, test.v, n)
	}
}

type conetest struct {
	p  Tuple4
	v  Tuple4
	t1 float64
	t2 float64
}
