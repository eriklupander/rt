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
	c.MinY = -5
	c.MaxY = 3
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
	c.MinY = -5
	c.MaxY = 3
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

func TestBoundingBox_MergeWith(t *testing.T) {
	b1 := NewBoundingBoxF(-5, -2, 0, 7, 4, 4)
	b2 := NewBoundingBoxF(8, -7, -2, 14, 2, 8)
	b1.MergeWith(b2)
	assert.Equal(t, NewPoint(-5, -7, -2), b1.Min)
	assert.Equal(t, NewPoint(14, 4, 8), b1.Max)
}

func TestBoundingBoxContainsPoint(t *testing.T) {

	BoundingBox := NewBoundingBoxF(5, -2, 0, 11, 4, 7)
	tests := []struct {
		point  Tuple4
		result bool
	}{
		{NewPoint(5, -2, 0), true},
		{NewPoint(11, 4, 7), true},
		{NewPoint(8, 1, 3), true},
		{NewPoint(3, 0, 3), false},
		{NewPoint(8, -4, 3), false},
		{NewPoint(8, 1, -1), false},
		{NewPoint(13, 1, 3), false},
		{NewPoint(8, 5, 3), false},
		{NewPoint(8, 1, 8), false},
	}

	for _, tc := range tests {
		res := BoundingBox.ContainsPoint(tc.point)
		assert.Equal(t, tc.result, res)
	}
}

func TestBoxContainsBox(t *testing.T) {

	BoundingBox := NewBoundingBoxF(5, -2, 0, 11, 4, 7)
	tests := []struct {
		min    Tuple4
		max    Tuple4
		result bool
	}{
		{NewPoint(5, -2, 0), NewPoint(11, 4, 7), true},
		{NewPoint(6, -1, 1), NewPoint(10, 3, 6), true},
		{NewPoint(4, -3, -1), NewPoint(10, 3, 6), false},
		{NewPoint(6, -1, 1), NewPoint(12, 5, 8), false},
	}

	for _, tc := range tests {
		res := BoundingBox.ContainsBox(NewBoundingBox(tc.min, tc.max))
		assert.Equal(t, tc.result, res)
	}
}

func TestTransformBoundingBox(t *testing.T) {
	box := NewBoundingBoxF(-1, -1, -1, 1, 1, 1)
	m1 := Multiply(RotateX(math.Pi/4), RotateY(math.Pi/4))
	box2 := TransformBoundingBox(box, m1)

	assert.InEpsilon(t, -1.4142, box2.Min[0], Epsilon)
	assert.InEpsilon(t, -1.7071, box2.Min[1], Epsilon)
	assert.InEpsilon(t, -1.7071, box2.Min[2], Epsilon)
	assert.InEpsilon(t, 1.4142, box2.Max[0], Epsilon)
	assert.InEpsilon(t, 1.7071, box2.Max[1], Epsilon)
	assert.InEpsilon(t, 1.7071, box2.Max[2], Epsilon)
}

func TestQueryBBTransformInParentSpace(t *testing.T) {
	shape := NewSphere()
	shape.SetTransform(Translate(1, -3, 5))
	shape.SetTransform(Scale(0.5, 2, 4))
	box := ParentSpaceBounds(shape)
	assert.Equal(t, NewPoint(0.5, -5, 1), box.Min)
	assert.Equal(t, NewPoint(1.5, -1, 9), box.Max)
}

func TestGroupBoundingBoxContainsAllItsChildren(t *testing.T) {

	s := NewSphere()
	s.SetTransform(Translate(2, 5, -3))
	s.SetTransform(Scale(2, 2, 2))

	c := NewCylinder()
	c.MinY = -2
	c.MaxY = 2
	c.SetTransform(Translate(-4, -1, 4))
	c.SetTransform(Scale(0.5, 1, 0.5))
	g := NewGroup()
	g.AddChild(s)
	g.AddChild(c)
	box := BoundsOf(g)
	assert.Equal(t, NewPoint(-4.5, -3, -5), box.Min)
	assert.Equal(t, NewPoint(4, 7, 4.5), box.Max)
}

func TestCSGBoundingBoxContainsAllItsChildren(t *testing.T) {

	left := NewSphere()
	right := NewSphere()
	right.SetTransform(Translate(2, 3, 4))
	csg := NewCSG("difference", left, right)
	box := BoundsOf(csg)
	assert.Equal(t, NewPoint(-1, -1, -1), box.Min)
	assert.Equal(t, NewPoint(3, 4, 5), box.Max)
}

func TestIntersectBoundingBoxWithRayAtOrigin(t *testing.T) {

	box := NewBoundingBoxF(-1, -1, -1, 1, 1, 1)

	testcases := []struct {
		origin    Tuple4
		direction Tuple4
		result    bool
	}{
		{NewPoint(5, 0.5, 0), NewVector(-1, 0, 0), true},
		{NewPoint(-5, 0.5, 0), NewVector(1, 0, 0), true},
		{NewPoint(0.5, 5, 0), NewVector(0, -1, 0), true},
		{NewPoint(0.5, -5, 0), NewVector(0, 1, 0), true},
		{NewPoint(0.5, 0, 5), NewVector(0, 0, -1), true},
		{NewPoint(0.5, 0, -5), NewVector(0, 0, 1), true},
		{NewPoint(0, 0.5, 0), NewVector(0, 0, 1), true},
		{NewPoint(-2, 0, 0), NewVector(2, 4, 6), false},
		{NewPoint(0, -2, 0), NewVector(6, 2, 4), false},
		{NewPoint(0, 0, -2), NewVector(4, 6, 2), false},
		{NewPoint(2, 0, 2), NewVector(0, 0, -1), false},
		{NewPoint(0, 2, 2), NewVector(0, -1, 0), false},
		{NewPoint(2, 2, 0), NewVector(-1, 0, 0), false},
	}

	for _, tc := range testcases {
		direction := Normalize(tc.direction)
		r := NewRay(tc.origin, direction)
		assert.Equal(t, tc.result, IntersectRayWithBox(r, box))
	}
}

func TestIntersectNonCubicBoundingBoxWithRay(t *testing.T) {

	box := NewBoundingBoxF(5, -2, 0, 11, 4, 7)

	testcases := []struct {
		origin    Tuple4
		direction Tuple4
		result    bool
	}{
		{NewPoint(15, 1, 2), NewVector(-1, 0, 0), true},
		{NewPoint(-5, -1, 4), NewVector(1, 0, 0), true},
		{NewPoint(7, 6, 5), NewVector(0, -1, 0), true},
		{NewPoint(9, -5, 6), NewVector(0, 1, 0), true},
		{NewPoint(8, 2, 12), NewVector(0, 0, -1), true},
		{NewPoint(6, 0, -5), NewVector(0, 0, 1), true},
		{NewPoint(8, 1, 3.5), NewVector(0, 0, 1), true},
		{NewPoint(9, -1, -8), NewVector(2, 4, 6), false},
		{NewPoint(8, 3, -4), NewVector(6, 2, 4), false},
		{NewPoint(9, -1, -2), NewVector(4, 6, 2), false},
		{NewPoint(4, 0, 9), NewVector(0, 0, -1), false},
		{NewPoint(8, 6, -1), NewVector(0, -1, 0), false},
		{NewPoint(12, 5, 4), NewVector(-1, 0, 0), false},
	}

	for _, tc := range testcases {
		direction := Normalize(tc.direction)
		r := NewRay(tc.origin, direction)
		assert.Equal(t, tc.result, IntersectRayWithBox(r, box))
	}
}

func TestIntersectRayGroupWithMiss(t *testing.T) {
	s := NewSphere()
	g := NewGroup()
	g.AddChild(s)
	g.Bounds()
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 1, 0))
	in := NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 0)) // Pass this as pointer for intermediate calc
	IntersectRayWithShapePtr(g, r, &in)

	// savedRay should have default values if the sphere's intersect was not called
	assert.Equal(t, 0.0, s.savedRay.Origin[0])
	assert.Equal(t, 0.0, s.savedRay.Origin[1])
	assert.Equal(t, 0.0, s.savedRay.Origin[2])
	assert.Equal(t, 0.0, s.savedRay.Direction[0])
	assert.Equal(t, 0.0, s.savedRay.Direction[1])
	assert.Equal(t, 0.0, s.savedRay.Direction[2])

}

func TestIntersectRayGroupWithHit(t *testing.T) {
	s := NewSphere()
	g := NewGroup()
	g.AddChild(s)
	g.Bounds()
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))
	in := NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 0)) // Pass this as pointer for intermediate calc
	IntersectRayWithShapePtr(g, r, &in)

	// savedRay should have val form ray if the sphere's intersect was called
	assert.Equal(t, 0.0, s.savedRay.Direction[0])
	assert.Equal(t, 0.0, s.savedRay.Direction[1])
	assert.Equal(t, 1.0, s.savedRay.Direction[2])

}

func TestIntersectRayWithCSGMissesBox(t *testing.T) {
	left := NewSphere()
	right := NewSphere()
	csg := NewCSG("difference", left, right)
	csg.Bounds()
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 1, 0))
	in := NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 0)) // Pass this as pointer for intermediate calc
	IntersectRayWithShapePtr(csg, r, &in)
	assert.Equal(t, 0.0, left.savedRay.Direction[0])
	assert.Equal(t, 0.0, right.savedRay.Direction[0])
	assert.Equal(t, 0.0, left.savedRay.Direction[1])
	assert.Equal(t, 0.0, right.savedRay.Direction[1])
	assert.Equal(t, 0.0, left.savedRay.Direction[2])
	assert.Equal(t, 0.0, right.savedRay.Direction[2])
}

func TestIntersectRayWithCSGHitsBox(t *testing.T) {
	left := NewSphere()
	right := NewSphere()
	csg := NewCSG("difference", left, right)
	csg.Bounds()
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))
	in := NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 0)) // Pass this as pointer for intermediate calc
	IntersectRayWithShapePtr(csg, r, &in)
	assert.Equal(t, 1.0, left.savedRay.Direction[2])
	assert.Equal(t, 1.0, right.savedRay.Direction[2])
}
