package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetup(t *testing.T) {
	st := DefaultSmoothTriangle()
	assert.Equal(t, NewPoint(0, 1, 0), st.P1)
}
func TestSmoothTriWithUV(t *testing.T) {
	st := DefaultSmoothTriangle()
	i := NewIntersectionUV(3.5, st, 0.2, 0.4)
	assert.Equal(t, 0.2, i.U)
	assert.Equal(t, 0.4, i.V)
}
func TestIntersectWithTriStoresUV(t *testing.T) {
	tri := DefaultSmoothTriangle()
	r := NewRay(NewPoint(-0.2, 0.3, -2), NewVector(0, 0, 1))
	xs := tri.IntersectLocal(r)
	assert.InEpsilon(t, 0.45, xs[0].U, Epsilon)
	assert.InEpsilon(t, 0.25, xs[0].V, Epsilon)
}
func TestInterpolatedNormal(t *testing.T) {
	tri := DefaultSmoothTriangle()
	i := NewIntersectionUV(1, tri, 0.45, 0.25)
	n := NormalAt(tri, NewPoint(0, 0, 0), &i)
	assert.InEpsilon(t, -0.5547, n.Get(0), Epsilon)
	assert.InEpsilon(t, 0.83205, n.Get(1), Epsilon)
	assert.Equal(t, 0.0, n.Get(2))
}
func TestPrepareNormalOnSmoothTri(t *testing.T) {
	tri := DefaultSmoothTriangle()
	i := NewIntersectionUV(1.0, tri, 0.45, 0.25)
	r := NewRay(NewPoint(-0.2, 0.3, -2), NewVector(0, 0, 1))
	xs := []Intersection{i}
	comps := PrepareComputationForIntersection(i, r, xs...)
	assert.InEpsilon(t, -0.5547, comps.NormalVec.Get(0), Epsilon)
	assert.InEpsilon(t, 0.83205, comps.NormalVec.Get(1), Epsilon)

	/*
		When i ← intersection_with_uv(1, tri, 0.45, 0.25)
		And r ← ray(point(-0.2, 0.3, -2), vector(0, 0, 1))
		And xs ← intersections(i)
		And comps ← prepare_computations(i, r, xs)
		Then comps.normalv = vector(-0.5547, 0.83205, 0)
	*/
}
