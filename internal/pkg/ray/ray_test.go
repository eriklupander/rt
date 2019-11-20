package ray

import (
	"github.com/eriklupander/rt/internal/pkg/mat"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRay(t *testing.T) {
	r := New(mat.NewPoint(1, 2, 3), mat.NewVector(4, 5, 6))
	assert.True(t, mat.TupleEquals(r.Origin, mat.NewPoint(1, 2, 3)))
	assert.True(t, mat.TupleEquals(r.Direction, mat.NewVector(4, 5, 6)))
}

func TestDistanceFromPoint(t *testing.T) {
	r := New(mat.NewPoint(2, 3, 4), mat.NewVector(1, 0, 0))
	p1 := Position(r, 0)
	assert.Equal(t, mat.NewPoint(2, 3, 4), p1)
}

func TestIntersectSphere(t *testing.T) {
	r := New(mat.NewPoint(0, 0, -5), mat.NewVector(0, 0, 1))
	s := mat.NewSphere()
	interects := IntersectRayWithSphere(s, r)
	assert.Len(t, interects, 2)
	assert.Equal(t, 4.0, interects[0].T)
	assert.Equal(t, 6.0, interects[1].T)
}

func TestIntersectSphereAtTangent(t *testing.T) {
	r := New(mat.NewPoint(0, 1, -5), mat.NewVector(0, 0, 1))
	s := mat.NewSphere()
	interects := IntersectRayWithSphere(s, r)
	assert.Len(t, interects, 2)
	assert.Equal(t, 5.0, interects[0].T)
	assert.Equal(t, 5.0, interects[1].T)
}

func TestIntersectMissSphere(t *testing.T) {
	r := New(mat.NewPoint(0, 2, -5), mat.NewVector(0, 0, 1))
	s := mat.NewSphere()
	interects := IntersectRayWithSphere(s, r)
	assert.Len(t, interects, 0)
}

func TestIntersectSphereWhenOriginatingFromCenterOfSphere(t *testing.T) {
	r := New(mat.NewPoint(0, 0, 0), mat.NewVector(0, 0, 1))
	s := mat.NewSphere()
	interects := IntersectRayWithSphere(s, r)
	assert.Len(t, interects, 2)
	assert.Equal(t, -1.0, interects[0].T)
	assert.Equal(t, 1.0, interects[1].T)
}

func TestIntersectSphereBehindRay(t *testing.T) {
	r := New(mat.NewPoint(0, 0, 5), mat.NewVector(0, 0, 1))
	s := mat.NewSphere()
	interects := IntersectRayWithSphere(s, r)
	assert.Len(t, interects, 2)
	assert.Equal(t, -6.0, interects[0].T)
	assert.Equal(t, -4.0, interects[1].T)
}

func TestHitWhenAllIntersectsHavePositiveT(t *testing.T) {
	s := mat.NewSphere()
	i1 := NewIntersection(1.0, s)
	i2 := NewIntersection(2.0, s)
	xs := []Intersection{*i1, *i2}
	i := Hit(xs)
	assert.True(t, IntersectionEqual(*i, *i1))
}

func TestHitWhenSomeIntersectsHaveNegativeT(t *testing.T) {
	s := mat.NewSphere()
	i1 := NewIntersection(-1.0, s)
	i2 := NewIntersection(1.0, s)
	xs := []Intersection{*i1, *i2}
	i := Hit(xs)
	assert.True(t, IntersectionEqual(*i, *i2))
}

func TestHitWhenAllIntersectsHaveNegativeT(t *testing.T) {
	s := mat.NewSphere()
	i1 := NewIntersection(-2.0, s)
	i2 := NewIntersection(-1.0, s)
	xs := []Intersection{*i1, *i2}
	i := Hit(xs)
	assert.Nil(t, i)
}

func TestHitIsLowestNonNegativeT(t *testing.T) {
	s := mat.NewSphere()
	i1 := NewIntersection(5.0, s)
	i2 := NewIntersection(7.0, s)
	i3 := NewIntersection(-3, s)
	i4 := NewIntersection(2.0, s)
	xs := []Intersection{*i1, *i2, *i3, *i4}
	i := Hit(xs)
	assert.True(t, IntersectionEqual(*i, *i4))
}

func TestTranslateRay(t *testing.T) {
	r := New(mat.NewPoint(1, 2, 3), mat.NewVector(0, 1, 0))
	m1 := mat.Translate(3, 4, 5)
	r2 := Transform(r, m1)
	assert.True(t, mat.TupleEquals(r2.Origin, mat.NewPoint(4, 6, 8)))
	assert.True(t, mat.TupleEquals(r2.Direction, mat.NewVector(0, 1, 0)))
}

func TestScaleRay(t *testing.T) {
	r := New(mat.NewPoint(1, 2, 3), mat.NewVector(0, 1, 0))
	m1 := mat.Scale(2, 3, 4)
	r2 := Transform(r, m1)
	assert.True(t, mat.TupleEquals(r2.Origin, mat.NewPoint(2, 6, 12)))
	assert.True(t, mat.TupleEquals(r2.Direction, mat.NewVector(0, 3, 0)))
}

func TestIntersectScaledSphereWithRay(t *testing.T) {
	r := New(mat.NewPoint(0, 0, -5), mat.NewVector(0, 0, 1))
	s := mat.NewSphere()
	mat.SetTransform(&s, mat.Scale(2, 2, 2))
	intersections := IntersectRayWithSphere(s, r)
	assert.Len(t, intersections, 2)
	assert.Equal(t, 3.0, intersections[0].T)
	assert.Equal(t, 7.0, intersections[1].T)
}

func TestIntersectTranslatedSphereWithRay(t *testing.T) {
	r := New(mat.NewPoint(0, 0, -5), mat.NewVector(0, 0, 1))
	s := mat.NewSphere()
	mat.SetTransform(&s, mat.Translate(5, 0, 0))
	intersections := IntersectRayWithSphere(s, r)
	assert.Len(t, intersections, 0)
}
