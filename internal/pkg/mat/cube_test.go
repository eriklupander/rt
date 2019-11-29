package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCubeRayHits(t *testing.T) {
	c := NewCube()

	tc := []cubetest{
		{name: "+x", p: NewPoint(5, 0.5, 0), v: NewVector(-1, 0, 0), t1: 4.0, t2: 6.0},
		{name: "-x", p: NewPoint(-5, 0.5, 0), v: NewVector(1, 0, 0), t1: 4.0, t2: 6.0},
		{name: "+y", p: NewPoint(0.5, 5, 0), v: NewVector(0, -1, 0), t1: 4.0, t2: 6.0},
		{name: "-y", p: NewPoint(0.5, -5, 0), v: NewVector(0, 1, 0), t1: 4.0, t2: 6.0},
		{name: "+z", p: NewPoint(0.5, 0, 5), v: NewVector(0, 0, -1), t1: 4.0, t2: 6.0},
		{name: "-z", p: NewPoint(0.5, 0, -5), v: NewVector(0, 0, 1), t1: 4.0, t2: 6.0},
	}
	for _, test := range tc {
		r := NewRay(test.p, test.v)
		xs := c.IntersectLocal(r)
		assert.Equal(t, test.t1, xs[0].T)
		assert.Equal(t, test.t2, xs[1].T)
	}
}

func TestCubeRayMisses(t *testing.T) {
	c := NewCube()

	tc := []cubetest{
		{p: NewPoint(-2, 0, 0), v: NewVector(0.2673, 0.5345, 0.8018)},
		{p: NewPoint(0, -2, 0), v: NewVector(0.8018, 0.2673, 0.5345)},
		{p: NewPoint(0, 0, -2), v: NewVector(0.5345, 0.8018, 0.2673)},
		{p: NewPoint(2, 0, 2), v: NewVector(0, 0, -1)},
		{p: NewPoint(0, 2, 2), v: NewVector(0, -1, 0)},
		{p: NewPoint(2, 2, 0), v: NewVector(-1, 0, 0)},
	}
	for _, test := range tc {
		r := NewRay(test.p, test.v)
		xs := c.IntersectLocal(r)
		assert.Len(t, xs, 0)
	}
}

func TestCubeNormal(t *testing.T) {
	c := NewCube()

	tc := []cubetest{
		{p: NewPoint(1, 0.5, -0.8), v: NewVector(1, 0, 0)},
		{p: NewPoint(-1, -0.2, 0.9), v: NewVector(-1, 0, 0)},
		{p: NewPoint(-0.4, 1, -0.1), v: NewVector(0, 1, 0)},
		{p: NewPoint(0.3, -1, -0.7), v: NewVector(0, -1, 0)},
		{p: NewPoint(-0.6, 0.3, 1), v: NewVector(0, 0, 1)},
		{p: NewPoint(0.4, 0.4, -1), v: NewVector(0, 0, -1)},
		{p: NewPoint(1, 1, 1), v: NewVector(1, 0, 0)},
		{p: NewPoint(-1, -1, -1), v: NewVector(-1, 0, 0)},
	}
	for _, test := range tc {
		n := c.NormalAtLocal(test.p)
		assert.Equal(t, test.v, n)
	}
}

type cubetest struct {
	name string
	p    Tuple4
	v    Tuple4
	t1   float64
	t2   float64
}
