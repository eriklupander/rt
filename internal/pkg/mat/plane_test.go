package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPlane(t *testing.T) {
	pl := NewPlane()
	n1 := pl.NormalAtLocal(NewPoint(0, 0, 0))
	n2 := pl.NormalAtLocal(NewPoint(10, 0, -10))
	n3 := pl.NormalAtLocal(NewPoint(-6, 0, 150))
	assert.True(t, TupleEquals(n1, NewVector(0, 1, 0)))
	assert.True(t, TupleEquals(n2, NewVector(0, 1, 0)))
	assert.True(t, TupleEquals(n3, NewVector(0, 1, 0)))

}

func TestPlane_IntersectLocalParallellMisses(t *testing.T) {
	pl := NewPlane()
	r := NewRay(NewPoint(0, 10, 0), NewVector(0, 0, 1))
	xs := pl.IntersectLocal(r)
	assert.Len(t, xs, 0)
}
func TestPlane_IntersectLocalCoplanarMisses(t *testing.T) {
	pl := NewPlane()
	r := NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 1))
	xs := pl.IntersectLocal(r)
	assert.Len(t, xs, 0)
}
func TestPlane_IntersectLocalFromAbove(t *testing.T) {
	pl := NewPlane()
	r := NewRay(NewPoint(0, 1, 0), NewVector(0, -1, 0))
	xs := pl.IntersectLocal(r)
	assert.Len(t, xs, 1)
	assert.Equal(t, 1.0, xs[0].T)
	assert.Equal(t, pl.ID(), xs[0].S.ID())
}
func TestPlane_IntersectLocalFromBelow(t *testing.T) {
	pl := NewPlane()
	r := NewRay(NewPoint(0, -1, 0), NewVector(0, 1, 0))
	xs := pl.IntersectLocal(r)
	assert.Len(t, xs, 1)
	assert.Equal(t, 1.0, xs[0].T)
	assert.Equal(t, pl.ID(), xs[0].S.ID())
}
