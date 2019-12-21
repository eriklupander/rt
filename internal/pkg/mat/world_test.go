package mat

import (
	"github.com/stretchr/testify/assert"
	"math"
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

// Page 95
func TestShadeIntersection(t *testing.T) {
	w := NewDefaultWorld()
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))
	i := Intersection{T: 4.0, S: w.Objects[0]}

	comps := PrepareComputationForIntersection(i, r)
	color := ShadeHit(w, comps, 1, 1)
	assert.InEpsilon(t, 0.38066, color.Get(0), Epsilon)
	assert.InEpsilon(t, 0.47583, color.Get(1), Epsilon)
	assert.InEpsilon(t, 0.2855, color.Get(2), Epsilon)
}

func TestShadeIntersectionFromInside(t *testing.T) {
	w := NewDefaultWorld()
	w.Light = []Light{NewLight(NewPoint(0, 0.25, 0), NewColor(1, 1, 1))}

	r := NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 1))

	i := Intersection{T: 0.5, S: w.Objects[1]}

	comps := PrepareComputationForIntersection(i, r)
	color := ShadeHit(w, comps, 1, 1)
	assert.InEpsilon(t, 0.90498, color.Get(0), Epsilon)
	assert.InEpsilon(t, 0.90498, color.Get(1), Epsilon)
	assert.InEpsilon(t, 0.90498, color.Get(2), Epsilon)
}

func TestColorWhenRayMiss(t *testing.T) {
	w := NewDefaultWorld()
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 1, 0))
	color := ColorAt(w, r, 1, 1)
	assert.Equal(t, color, NewColor(0, 0, 0))
}

func TestColorWhenRayHits(t *testing.T) {
	w := NewDefaultWorld()
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))
	color := ColorAt(w, r, 1, 1)
	assert.InEpsilon(t, 0.38066, color.Get(0), Epsilon)
	assert.InEpsilon(t, 0.47583, color.Get(1), Epsilon)
	assert.InEpsilon(t, 0.2855, color.Get(2), Epsilon)
}

// Page 97
func TestColorWhenCastWithinSphereAtInsideSphere(t *testing.T) {
	w := NewDefaultWorld()
	w.Objects[0].SetMaterial(NewMaterial(NewColor(0.8, 1.0, 0.6), 1.0, 0.7, 0.2, 200))
	w.Objects[1].SetMaterial(NewMaterial(NewColor(0.8, 1.0, 0.6), 1.0, 0.7, 0.2, 200))

	r := NewRay(NewPoint(0, 0, 0.75), NewVector(0, 0, -1))
	color := ColorAt(w, r, 1, 1)
	assert.InEpsilon(t, w.Objects[1].GetMaterial().Color.Get(0), color.Get(0), Epsilon)
	assert.InEpsilon(t, w.Objects[1].GetMaterial().Color.Get(1), color.Get(1), Epsilon)
	assert.InEpsilon(t, w.Objects[1].GetMaterial().Color.Get(2), color.Get(2), Epsilon)
}

func TestPointNotInShadow(t *testing.T) {
	w := NewDefaultWorld()
	p := NewPoint(0, 10, 10)
	assert.False(t, PointInShadow(w, w.Light[0], p))
}
func TestPointInShadow(t *testing.T) {
	w := NewDefaultWorld()
	p := NewPoint(10, -10, 10)
	assert.True(t, PointInShadow(w, w.Light[0], p))
}
func TestPointNotInShadowWhenBehindLight(t *testing.T) {
	w := NewDefaultWorld()
	p := NewPoint(-20, 20, -20)
	assert.False(t, PointInShadow(w, w.Light[0], p))
}
func TestPointNotInShadowWhenBehindPoint(t *testing.T) {
	w := NewDefaultWorld()
	p := NewPoint(-2, 2, -2)
	assert.False(t, PointInShadow(w, w.Light[0], p))
}

// Big one on page 114
func TestWorldWithShadowTest(t *testing.T) {
	w := NewDefaultWorld()
	w.Light = []Light{NewLight(NewPoint(0, 0, -10), NewColor(1, 1, 1))}
	s := NewSphere()
	w.Objects = append(w.Objects, s)
	s2 := NewSphere()
	s2.Transform = Multiply(s2.Transform, Translate(0, 0, 10))
	w.Objects = append(w.Objects, s2)

	r := NewRay(NewPoint(0, 0, 5), NewVector(0, 0, 1))
	i := NewIntersection(4, s2)
	comps := PrepareComputationForIntersection(i, r)
	color := ShadeHit(w, comps, 1, 1)
	color.Elems[3] = 1 // just a fix for me using Tuple4 to represent colors...
	assert.Equal(t, NewColor(0.1, 0.1, 0.1), color)
}

func TestHitOffsetToFixAcne(t *testing.T) {
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))
	s := NewSphere()
	s.Transform = Multiply(s.Transform, Translate(0, 0, 1))
	i := NewIntersection(5, s)
	comps := PrepareComputationForIntersection(i, r)
	assert.True(t, comps.OverPoint.Get(2) < -Epsilon/2)
	assert.True(t, comps.Point.Get(2) > comps.OverPoint.Get(2))
}

func TestOpaqueRefraction(t *testing.T) {
	w := NewDefaultWorld()
	s1 := w.Objects[0]
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))
	xs := []Intersection{
		NewIntersection(4, s1), NewIntersection(6, s1),
	}
	comps := PrepareComputationForIntersection(xs[0], r, xs...)
	color := RefractedColor(w, comps, 5)
	assert.Equal(t, black, color)
}

func TestRefractiveColorAndMaxRecursionDepth(t *testing.T) {
	w := NewDefaultWorld()
	s1 := w.Objects[0]
	material := NewDefaultMaterial()
	material.Transparency = 1.0
	material.RefractiveIndex = 1.5
	s1.SetMaterial(material)

	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))
	xs := []Intersection{
		NewIntersection(4, s1), NewIntersection(6, s1),
	}
	comps := PrepareComputationForIntersection(xs[0], r, xs...)
	color := RefractedColor(w, comps, 0)
	assert.Equal(t, black, color)
}

func TestTotalInternalRefraction(t *testing.T) {
	w := NewDefaultWorld()
	s1 := w.Objects[0]
	material := NewDefaultMaterial()
	material.Transparency = 1.0
	material.RefractiveIndex = 1.5
	s1.SetMaterial(material)
	r := NewRay(NewPoint(0, 0, math.Sqrt(2)/2), NewVector(0, 1, 0))

	xs := []Intersection{
		NewIntersection(-math.Sqrt(2)/2, s1), NewIntersection(math.Sqrt(2)/2, s1),
	}

	comps := PrepareComputationForIntersection(xs[1], r, xs...)
	color := RefractedColor(w, comps, 5)
	assert.Equal(t, black, color)
}

func TestRefractedColorWithRefractedRay(t *testing.T) {
	w := NewDefaultWorld()
	s1 := w.Objects[0]
	material := NewDefaultMaterial()
	material.Ambient = 1.0
	material.Pattern = NewTestPattern()
	s1.SetMaterial(material)

	s2 := w.Objects[1]
	material2 := NewDefaultMaterial()
	material2.Transparency = 1.0
	material2.RefractiveIndex = 1.5
	s2.SetMaterial(material2)

	r := NewRay(NewPoint(0, 0, 0.1), NewVector(0, 1, 0))

	xs := []Intersection{
		NewIntersection(-0.9899, s1),
		NewIntersection(-0.4899, s2),
		NewIntersection(0.4899, s2),
		NewIntersection(0.4899, s1),
	}

	comps := PrepareComputationForIntersection(xs[2], r, xs...)
	color := RefractedColor(w, comps, 5)
	assert.Equal(t, 0.0, color.Get(0))
	assert.InEpsilon(t, 0.99888, color.Get(1), Epsilon)
	assert.InEpsilon(t, 0.04725, color.Get(2), Epsilon*10)
}

func TestShadeHitWithRefractedMaterial(t *testing.T) {
	w := NewDefaultWorld()
	floor := NewPlane()
	floor.SetTransform(Translate(0, -1, 0))
	mat1 := NewDefaultMaterial()
	mat1.Transparency = 0.5
	mat1.RefractiveIndex = 1.5
	floor.SetMaterial(mat1)
	w.Objects = append(w.Objects, floor)

	ball := NewSphere()
	mat2 := NewDefaultMaterial()
	mat2.Color = NewColor(1, 0, 0)
	mat2.Ambient = 0.5
	ball.SetMaterial(mat2)
	ball.SetTransform(Translate(0, -3.5, -0.5))
	w.Objects = append(w.Objects, ball)

	ray := NewRay(NewPoint(0, 0, -3), NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	xs := []Intersection{
		NewIntersection(math.Sqrt(2), floor),
	}
	comps := PrepareComputationForIntersection(xs[0], ray, xs...)
	color := ShadeHit(w, comps, 5, 5)
	assert.InEpsilon(t, 0.93642, color.Get(0), Epsilon)
	assert.InEpsilon(t, 0.68642, color.Get(1), Epsilon)
	assert.InEpsilon(t, 0.68642, color.Get(2), Epsilon)
}

func TestShadeHitWhenBothTransparentAndRefractive(t *testing.T) {
	w := NewDefaultWorld()
	r := NewRay(NewPoint(0, 0, -3), NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))

	floor := NewPlane()
	floor.SetTransform(Translate(0, -1, 0))
	mat1 := NewDefaultMaterial()
	mat1.Reflectivity = 0.5
	mat1.Transparency = 0.5
	mat1.RefractiveIndex = 1.5
	floor.SetMaterial(mat1)
	w.Objects = append(w.Objects, floor)

	ball := NewSphere()
	ball.SetTransform(Translate(0, -3.5, -0.5))
	mat2 := NewDefaultMaterial()
	mat2.Color = NewColor(1, 0, 0)
	mat2.Ambient = 0.5
	ball.SetMaterial(mat2)
	w.Objects = append(w.Objects, ball)

	xs := []Intersection{
		NewIntersection(math.Sqrt(2), floor),
	}
	color := ShadeHit(w, PrepareComputationForIntersection(xs[0], r, xs...), 5, 5)
	assert.InEpsilon(t, 0.93642, color.Get(0), Epsilon*3)
	assert.InEpsilon(t, 0.69643, color.Get(1), Epsilon)
	assert.InEpsilon(t, 0.69243, color.Get(2), Epsilon)
}
