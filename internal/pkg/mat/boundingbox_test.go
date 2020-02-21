package mat

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestNewEmptyBoundingBox(t *testing.T) {
	box := NewEmptyBoundingBox()
	assert.Equal(t, NewTupleOf(math.Inf(1), math.Inf(1), math.Inf(1), 1), box.Min)
	assert.Equal(t, NewTupleOf(math.Inf(-1), math.Inf(-1), math.Inf(-1), 1), box.Max)
}

func TestNewBoundingBoxWithVolume(t *testing.T) {
	box := NewBoundingBoxF(-1, -2, -3, 3, 2, 1)
	assert.Equal(t, NewPoint(-1, -2, -3), box.Min)
	assert.Equal(t, NewPoint(3, 2, 1), box.Max)
}

func TestAddPointToBoundingBox(t *testing.T) {

	box := NewEmptyBoundingBox()
	p1 := NewPoint(-5, 2, 0)
	p2 := NewPoint(7, 0, -3)
	box.Add(p1)
	box.Add(p2)
	assert.Equal(t, NewPoint(-5, 0, -3), box.Min)
	assert.Equal(t, NewPoint(7, 2, 0), box.Max)
}

func TestBoundsOfSphere(t *testing.T) {
	s := NewSphere()
	box := BoundsOf(s)
	assert.Equal(t, NewPoint(-1, -1, -1), box.Min)
	assert.Equal(t, NewPoint(1, 1, 1), box.Max)
}

func TestBoundsOfPlane(t *testing.T) {
	p := NewPlane()
	box := BoundsOf(p)
	assert.Equal(t, NewPoint(math.Inf(-1), 0, math.Inf(-1)), box.Min)
	assert.Equal(t, NewPoint(math.Inf(1), 0, math.Inf(1)), box.Max)
}

func TestBoundsOfCube(t *testing.T) {
	c := NewCube()
	box := BoundsOf(c)
	assert.Equal(t, NewPoint(-1, -1, -1), box.Min)
	assert.Equal(t, NewPoint(1, 1, 1), box.Max)
}

func TestBoundsOfInfiniteCylinder(t *testing.T) {
	c := NewCylinder()
	box := BoundsOf(c)
	assert.Equal(t, NewPoint(-1, math.Inf(-1), -1), box.Min)
	assert.Equal(t, NewPoint(1, math.Inf(1), 1), box.Max)
}

func TestBoundsOfFiniteCylinder(t *testing.T) {
	c := NewCylinder()
	c.minY = -5
	c.maxY = 3
	box := BoundsOf(c)
	assert.Equal(t, NewPoint(-1, -5, -1), box.Min)
	assert.Equal(t, NewPoint(1, 3, 1), box.Max)
}

func TestBoundsOfInfiniteCone(t *testing.T) {
	c := NewCone()
	box := BoundsOf(c)
	assert.Equal(t, NewPoint(math.Inf(-1), math.Inf(-1), math.Inf(-1)), box.Min)
	assert.Equal(t, NewPoint(math.Inf(1), math.Inf(1), math.Inf(1)), box.Max)
}

func TestBoundsOfFiniteCone(t *testing.T) {
	c := NewCone()
	c.minY = -5
	c.maxY = 3
	box := BoundsOf(c)
	assert.Equal(t, NewPoint(-5, -5, -5), box.Min)
	assert.Equal(t, NewPoint(5, 3, 5), box.Max)
}

func TestBoundsOfTriangle(t *testing.T) {
	p1 := NewPoint(-3, 7, 2)
	p2 := NewPoint(6, 2, -4)
	p3 := NewPoint(2, -1, -1)
	tri := NewTriangle(p1, p2, p3)
	box := BoundsOf(tri)
	assert.Equal(t, NewPoint(-3, -1, -4), box.Min)
	assert.Equal(t, NewPoint(6, 7, 2), box.Max)
}
