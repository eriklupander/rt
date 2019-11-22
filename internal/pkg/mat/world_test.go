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

// Page 95
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

// Page 97
func TestColorWhenCastWithinSphereAtInsideSphere(t *testing.T) {
	w := NewDefaultWorld()
	w.Objects[0].Material.Ambient = 1.0
	w.Objects[1].Material = NewMaterial(NewColor(0.8, 1.0, 0.6), 0.1, 0.7, 0.2, 200)
	w.Objects[1].Material.Ambient = 1.0

	r := NewRay(NewPoint(0, 0, 0.75), NewVector(0, 0, -1))
	color := ColorAt(w, r)
	assert.InEpsilon(t, w.Objects[1].Material.Color.Get(0), color.Get(0), Epsilon)
	assert.InEpsilon(t, w.Objects[1].Material.Color.Get(1), color.Get(1), Epsilon)
	assert.InEpsilon(t, w.Objects[1].Material.Color.Get(2), color.Get(2), Epsilon)
}

func TestPointNotInShadow(t *testing.T) {
	w := NewDefaultWorld()
	p := NewPoint(0, 10, 10)
	assert.False(t, PointInShadow(w, p))
}
func TestPointInShadow(t *testing.T) {
	w := NewDefaultWorld()
	p := NewPoint(10, -10, 10)
	assert.True(t, PointInShadow(w, p))
}
func TestPointNotInShadowWhenBehindLight(t *testing.T) {
	w := NewDefaultWorld()
	p := NewPoint(-20, 20, -20)
	assert.False(t, PointInShadow(w, p))
}
func TestPointNotInShadowWhenBehindPoint(t *testing.T) {
	w := NewDefaultWorld()
	p := NewPoint(-2, 2, -2)
	assert.False(t, PointInShadow(w, p))
}

// Big one on page 114
func TestWorldWithShadowTest(t *testing.T) {
	w := NewDefaultWorld()
	w.Light = NewLight(NewPoint(0, 0, -10), NewColor(1, 1, 1))
	s := NewSphere()
	w.Objects = append(w.Objects, s)
	s2 := NewSphere()
	s2.Transform = Multiply(s2.Transform, Translate(0, 0, 10))
	w.Objects = append(w.Objects, s2)

	r := NewRay(NewPoint(0, 0, 5), NewVector(0, 0, 1))
	i := NewIntersection(4, s2)
	comps := PrepareComputationForIntersection(i, r)
	color := ShadeHit(w, comps)
	assert.Equal(t, NewTuple4([]float64{0.1, 0.1, 0.1, 0.1}), color)
}

func TestHitOffsetToFixAcne(t *testing.T) {
	/*
		Given r ← ray(point(0, 0, -5), vector(0, 0, 1))
		And shape ← sphere() with:
		| transform | translation(0, 0, 1) |
		And i ← intersection(5, shape)
		When comps ← prepare_computations(i, r)
		Then comps.over_point.z < -EPSILON/2
		And comps.point.z > comps.over_point.z
	*/
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))
	s := NewSphere()
	s.Transform = Multiply(s.Transform, Translate(0, 0, 1))
	i := NewIntersection(5, s)
	comps := PrepareComputationForIntersection(i, r)
	assert.True(t, comps.OverPoint.Get(2) < -Epsilon/2)
	assert.True(t, comps.Point.Get(2) > comps.OverPoint.Get(2))
}
