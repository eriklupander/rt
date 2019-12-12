package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCylinderRayMisses(t *testing.T) {
	c := NewCylinder()

	tc := []cyltest{
		{p: NewPoint(1, 0, 0), v: NewVector(0, 1, 0)},
		{p: NewPoint(0, 0, 0), v: NewVector(0, 1, 0)},
		{p: NewPoint(0, 0, -5), v: NewVector(1, 1, 1)},
	}
	for _, test := range tc {
		r := NewRay(test.p, test.v)
		xs := c.IntersectLocal(r)
		assert.Len(t, xs, 0)
	}
}
func TestCylinderRayHits(t *testing.T) {
	c := NewCylinder()

	tc := []cyltest{
		{p: NewPoint(1, 0, -5), v: Normalize(NewVector(0, 0, 1)), t1: 5, t2: 5},
		{p: NewPoint(0, 0, -5), v: Normalize(NewVector(0, 0, 1)), t1: 4, t2: 6},
		{p: NewPoint(0.5, 0, -5), v: Normalize(NewVector(0.1, 1, 1)), t1: 6.80798191702732, t2: 7.088723439378861},
	}
	for _, test := range tc {
		r := NewRay(test.p, test.v)
		xs := c.IntersectLocal(r)
		assert.Len(t, xs, 2)
		assert.Equal(t, test.t1, xs[0].T)
		assert.Equal(t, test.t2, xs[1].T)
	}
}

func TestCylinderLocalNormal(t *testing.T) {
	c := NewCylinder()

	tc := []cyltest{
		{p: NewPoint(1, 0, 0), v: Normalize(NewVector(1, 0, 0))},
		{p: NewPoint(0, 5, -1), v: Normalize(NewVector(0, 0, -1))},
		{p: NewPoint(0, -2, 1), v: Normalize(NewVector(0, 0, 1))},
		{p: NewPoint(-1, 1, 0), v: Normalize(NewVector(-1, 0, 0))},
	}
	for _, test := range tc {
		n := c.NormalAtLocal(test.p, nil)
		assert.Equal(t, test.v, n)
	}
}

func TestIntersectCappedOpenCylinder(t *testing.T) {
	c := NewCylinderMM(1, 2)

	tc := []cyltest{
		{p: NewPoint(0, 1.5, 0), v: NewVector(0.1, 1, 0), t1: 0},
		{p: NewPoint(0, 3, -5), v: NewVector(0, 0, 1), t1: 0},
		{p: NewPoint(0, 0, -5), v: NewVector(0, 0, 1), t1: 0},
		{p: NewPoint(0, 2, -5), v: NewVector(0, 0, 1), t1: 0},
		{p: NewPoint(0, 1, -5), v: NewVector(0, 0, 1), t1: 0},
		{p: NewPoint(0, 1.5, -2), v: NewVector(0, 0, 1), t1: 2},
	}
	for _, test := range tc {
		xs := c.IntersectLocal(NewRay(test.p, Normalize(test.v)))
		assert.Len(t, xs, int(test.t1))
	}
}

func TestIntersectCappedClosedCylinder(t *testing.T) {
	c := NewCylinderMMC(1, 2, true)

	tc := []cyltest{
		{p: NewPoint(0, 3, 0), v: NewVector(0, -1, 0), t1: 2},
		{p: NewPoint(0, 3, -2), v: NewVector(0, -1, 2), t1: 2},
		{p: NewPoint(0, 4, -2), v: NewVector(0, -1, 1), t1: 2},
		{p: NewPoint(0, 0, -2), v: NewVector(0, 1, 2), t1: 2},
		{p: NewPoint(0, -1, -2), v: NewVector(0, 1, 1), t1: 2},
	}
	for _, test := range tc {
		xs := c.IntersectLocal(NewRay(test.p, Normalize(test.v)))
		assert.Len(t, xs, int(test.t1))
	}
}

func TestCylinderNormalAtCap(t *testing.T) {
	cyl := NewCylinderMMC(1, 2, true)

	tc := []cyltest{
		{p: NewPoint(0, 1, 0), v: NewVector(0, -1, 0), t1: 2},
		{p: NewPoint(0.5, 1, 0), v: NewVector(0, -1, 0), t1: 2},
		{p: NewPoint(0, 1, 0.5), v: NewVector(0, -1, 0), t1: 2},
		{p: NewPoint(0, 2, 0), v: NewVector(0, 1, 0), t1: 2},
		{p: NewPoint(0.5, 2, 0), v: NewVector(0, 1, 0), t1: 2},
		{p: NewPoint(0, 2, 0.5), v: NewVector(0, 1, 0), t1: 2},
	}

	for _, test := range tc {
		n := cyl.NormalAtLocal(test.p, nil)
		assert.Equal(t, test.v, n)
	}
}

type cyltest struct {
	p  Tuple4
	v  Tuple4
	t1 float64
	t2 float64
}
