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

func TestIntersectionAllowedUnionCSG(t *testing.T) {
	tc := []struct {
		op     string
		lhit   bool
		inl    bool
		inr    bool
		result bool
	}{
		{"union", true, true, true, false},
		{"union", true, true, false, true},
		{"union", true, false, true, false},
		{"union", true, false, false, true},
		{"union", false, true, true, false},
		{"union", false, true, false, false},
		{"union", false, false, true, true},
		{"union", false, false, false, true},
	}

	for _, test := range tc {
		assert.True(t, IntersectionAllowed(test.op, test.lhit, test.inl, test.inr) == test.result)
	}
}

func TestIntersectionAllowedIntersectionCSG(t *testing.T) {
	tc := []struct {
		op     string
		lhit   bool
		inl    bool
		inr    bool
		result bool
	}{
		{"intersection", true, true, true, true},
		{"intersection", true, true, false, false},
		{"intersection", true, false, true, true},
		{"intersection", true, false, false, false},
		{"intersection", false, true, true, true},
		{"intersection", false, true, false, true},
		{"intersection", false, false, true, false},
		{"intersection", false, false, false, false},
	}

	for _, test := range tc {
		assert.True(t, IntersectionAllowed(test.op, test.lhit, test.inl, test.inr) == test.result)
	}
}

func TestIntersectionAllowedDifferenceCSG(t *testing.T) {
	tc := []struct {
		op     string
		lhit   bool
		inl    bool
		inr    bool
		result bool
	}{
		{"difference", true, true, true, false},
		{"difference", true, true, false, true},
		{"difference", true, false, true, false},
		{"difference", true, false, false, true},
		{"difference", false, true, true, true},
		{"difference", false, true, false, true},
		{"difference", false, false, true, false},
		{"difference", false, false, false, false},
	}

	for _, test := range tc {
		assert.True(t, IntersectionAllowed(test.op, test.lhit, test.inl, test.inr) == test.result)
	}
}

func TestFilterCSGIntersections(t *testing.T) {
	s1 := NewSphere()
	s2 := NewCube()
	xs := []Intersection{
		NewIntersection(1, s1),
		NewIntersection(2, s2),
		NewIntersection(3, s1),
		NewIntersection(4, s2),
	}
	c1 := NewCSG("union", s1, s2)
	c2 := NewCSG("intersection", s1, s2)
	c3 := NewCSG("difference", s1, s2)

	r1 := FilterIntersections(c1, xs)
	r2 := FilterIntersections(c2, xs)
	r3 := FilterIntersections(c3, xs)

	assert.Len(t, r1, 2)
	assert.Len(t, r2, 2)
	assert.Len(t, r3, 2)

	assert.Equal(t, r1[0], xs[0])
	assert.Equal(t, r1[1], xs[3])
	assert.Equal(t, r2[0], xs[1])
	assert.Equal(t, r2[1], xs[2])
	assert.Equal(t, r3[0], xs[0])
	assert.Equal(t, r3[1], xs[1])
}
