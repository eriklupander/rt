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
	assert.True(t, TupleEquals(w.Light.Position, NewPoint(-10, 10, -10)))
	assert.True(t, TupleEquals(w.Light.Intensity, NewPoint(1, 1, 1)))
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
	assert.Equal(t, comps.Object.Id, s.Id)
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
	assert.Equal(t, comps.Object.Id, s.Id)
	assert.Equal(t, comps.Point, NewPoint(0, 0, 1))
	assert.Equal(t, comps.EyeVec, NewVector(0, 0, -1))
	assert.Equal(t, comps.NormalVec, NewVector(0, 0, -1))
	assert.True(t, comps.Inside)

}

func TestShadeIntersection(t *testing.T) {
	w := NewDefaultWorld()
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))
	i := Intersection{T: 4.0, S: w.Objects[0]}

	comps := PrepareComputationForIntersection(i, r)
	color := ShadeHit(w, comps)
	assert.InEpsilon(t, 0.38066, color.Get(0), Epsilon)
	assert.InEpsilon(t, 0.47583, color.Get(1), Epsilon)
	assert.InEpsilon(t, 0.2855, color.Get(2), Epsilon)
}

func TestShadeIntersectionFromInside(t *testing.T) {
	w := NewDefaultWorld()
	w.Light = NewLight(NewPoint(0, 0.25, 0), NewColor(1, 1, 1))

	r := NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 1))

	i := Intersection{T: 0.5, S: w.Objects[1]}

	comps := PrepareComputationForIntersection(i, r)
	color := ShadeHit(w, comps)
	assert.InEpsilon(t, 0.90498, color.Get(0), Epsilon)
	assert.InEpsilon(t, 0.90498, color.Get(1), Epsilon)
	assert.InEpsilon(t, 0.90498, color.Get(2), Epsilon)
}

func TestColorWhenRayMiss(t *testing.T) {
	w := NewDefaultWorld()
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 1, 0))
	color := ColorAt(w, r)
	assert.Equal(t, color, NewColor(0, 0, 0))
}

func TestColorWhenRayHits(t *testing.T) {
	w := NewDefaultWorld()
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))
	color := ColorAt(w, r)
	assert.InEpsilon(t, 0.38066, color.Get(0), Epsilon)
	assert.InEpsilon(t, 0.47583, color.Get(1), Epsilon)
	assert.InEpsilon(t, 0.2855, color.Get(2), Epsilon)
}

func TestColorWhenCastWithinSphereAtInsideSphere(t *testing.T) {
	w := NewDefaultWorld()
	w.Objects[0].Material.Ambient = 1.0
	w.Objects[1].Material.Ambient = 1.0

	r := NewRay(NewPoint(0, 0, 0.75), NewVector(0, 0, -1))
	color := ColorAt(w, r)
	assert.InEpsilon(t, w.Objects[1].Material.Color.Get(0), color.Get(0), Epsilon)
	assert.InEpsilon(t, w.Objects[1].Material.Color.Get(1), color.Get(1), Epsilon)
	assert.InEpsilon(t, w.Objects[1].Material.Color.Get(2), color.Get(2), Epsilon)
}
