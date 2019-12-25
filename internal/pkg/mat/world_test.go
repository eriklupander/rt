package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWorld(t *testing.T) {
	w := NewWorld()
	assert.NotNil(t, w)
}

func TestDefaultWorld(t *testing.T) {
	w := NewDefaultWorld()
	assert.Len(t, w.Objects, 2)
	assert.True(t, TupleEquals(w.Light[0].Position, NewPoint(-10, 10, -10)))
	assert.True(t, TupleEquals(w.Light[0].Intensity, NewPoint(1, 1, 1)))
}

func TestIntersectWorldWithRay(t *testing.T) {
	w := NewDefaultWorld()

	// this ray goes from 5 units in front of origo and points directly at the origo.
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))

	intersections := IntersectWithWorld(w, r)
	assert.Equal(t, intersections[0].T, 4.0)
	assert.Equal(t, intersections[1].T, 4.5)
	assert.Equal(t, intersections[2].T, 5.5)
	assert.Equal(t, intersections[3].T, 6.0)
}

func TestPrepareCompute(t *testing.T) {
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))
	s := NewSphere()
	xs := Intersection{T: 4.0, S: s}

	comps := PrepareComputationForIntersection(xs, r)
	assert.Equal(t, comps.T, xs.T)
	assert.Equal(t, comps.Object.ID(), s.Id)
	assert.Equal(t, comps.Point, NewPoint(0, 0, -1))
	assert.Equal(t, comps.EyeVec, NewVector(0, 0, -1))
	assert.Equal(t, comps.NormalVec, NewVector(0, 0, -1))
	assert.False(t, comps.Inside)

}

func TestPrepareComputeInsideHit(t *testing.T) {
	r := NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 1))
	s := NewSphere()
	xs := Intersection{T: 1.0, S: s}

	comps := PrepareComputationForIntersection(xs, r)
	assert.Equal(t, comps.T, xs.T)
	assert.Equal(t, comps.Object.ID(), s.Id)
	assert.Equal(t, comps.Point, NewPoint(0, 0, 1))
	assert.Equal(t, comps.EyeVec, NewVector(0, 0, -1))
	assert.Equal(t, comps.NormalVec, NewVector(0, 0, -1))
	assert.True(t, comps.Inside)

}

func TestHitOffsetToFixAcne(t *testing.T) {
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))
	s := NewSphere()
	s.SetTransform(Multiply(s.Transform, Translate(0, 0, 1)))
	i := NewIntersection(5, s)
	comps := PrepareComputationForIntersection(i, r)
	assert.True(t, comps.OverPoint.Get(2) < -Epsilon/2)
	assert.True(t, comps.Point.Get(2) > comps.OverPoint.Get(2))
}
