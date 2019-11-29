package mat

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestNewIntersection(t *testing.T) {
	sp := NewSphere()
	intersection := NewIntersection(3.5, sp)
	assert.Equal(t, 3.5, intersection.T)
	assert.Equal(t, sp.Id, intersection.S.ID())
}

func TestFindIntersectionsInDifferentMaterials(t *testing.T) {
	s1 := NewSphere()
	s1.SetMaterial(NewGlassMaterial(1.5))
	s1.SetTransform(Scale(2, 2, 2))

	s2 := NewSphere()
	s2.SetMaterial(NewGlassMaterial(2.0))
	s2.SetTransform(Translate(0, 0, -0.25))

	s3 := NewSphere()
	s3.SetMaterial(NewGlassMaterial(2.5))
	s3.SetTransform(Translate(0, 0, 0.25))

	ray := NewRay(NewPoint(0, 0, -4), NewVector(0, 0, 1))
	xs := []Intersection{
		{T: 2.0, S: s1},
		{T: 2.75, S: s2},
		{T: 3.25, S: s3},
		{T: 4.75, S: s2},
		{T: 5.25, S: s3},
		{T: 6.0, S: s1},
	}

	tc := []xsTestCase{
		{n1: 1.0, n2: 1.5},
		{n1: 1.5, n2: 2.0},
		{n1: 2.0, n2: 2.5},
		{n1: 2.5, n2: 2.5},
		{n1: 2.5, n2: 1.5},
		{n1: 1.5, n2: 1.0},
	}

	for i := 0; i < len(xs); i++ {
		comps := PrepareComputationForIntersection(xs[i], ray, xs...)
		assert.Equal(t, tc[i].n1, comps.N1, "wrong N1 at index %d", i)
		assert.Equal(t, tc[i].n2, comps.N2, "wrong N2 at index %d", i)
	}
}

type xsTestCase struct {
	n1 float64
	n2 float64
}

func TestIntesectionUnderPoint(t *testing.T) {
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))
	shape := NewSphere()
	shape.SetTransform(Translate(0, 0, 1))
	shape.SetMaterial(NewGlassMaterial(1.0))
	i := NewIntersection(5, shape)
	xs := []Intersection{i}
	comps := PrepareComputationForIntersection(i, r, xs...)
	assert.True(t, comps.UnderPoint.Get(2) > Epsilon/2)
	assert.True(t, comps.Point.Get(2) < comps.UnderPoint.Get(2))
}

func TestSchlickUnderTotalInternalReflection(t *testing.T) {
	shape := NewGlassSphere()
	r := NewRay(NewPoint(0, 0, math.Sqrt(2)/2), NewVector(0, 1, 0))
	xs := []Intersection{
		NewIntersection(-math.Sqrt(2)/2, shape),
		NewIntersection(math.Sqrt(2)/2, shape),
	}
	comps := PrepareComputationForIntersection(xs[1], r, xs...)
	reflectance := Schlick(comps)
	assert.Equal(t, 1.0, reflectance)
}

func TestSchlickWhenPerpendicular(t *testing.T) {
	/*
		The Schlick approximation with a perpendicular viewing angle
		Given shape ← glass_sphere()
		And r ← ray(point(0, 0, 0), vector(0, 1, 0))
		And xs ← intersections(-1:shape, 1:shape)
		When comps ← prepare_computations(xs[1], r, xs)
		And reflectance ← schlick(comps)
		Then reflectance = 0.04
	*/
	shape := NewGlassSphere()
	r := NewRay(NewPoint(0, 0, 0), NewVector(0, 1, 0))
	xs := []Intersection{
		NewIntersection(-1, shape),
		NewIntersection(1, shape),
	}
	comps := PrepareComputationForIntersection(xs[1], r, xs...)
	reflectance := Schlick(comps)
	assert.InEpsilon(t, 0.04, reflectance, Epsilon)
}

func TestSchlickWhenAngleIsSmall(t *testing.T) {
	/*
		Scenario: The Schlick approximation with small angle and n2 > n1
		Given shape ← glass_sphere()
		And r ← ray(point(0, 0.99, -2), vector(0, 0, 1))
		And xs ← intersections(1.8589:shape)
		When comps ← prepare_computations(xs[0], r, xs)
		And reflectance ← schlick(comps)
		Then reflectance = 0.48873
	*/
	shape := NewGlassSphere()
	r := NewRay(NewPoint(0, 0.99, -2), NewVector(0, 0, 1))
	xs := []Intersection{
		NewIntersection(1.8589, shape),
	}
	comps := PrepareComputationForIntersection(xs[0], r, xs...)
	reflectance := Schlick(comps)
	assert.Equal(t, 0.4887308101221217, reflectance)
}
