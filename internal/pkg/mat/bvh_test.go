package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitPerfectCube(t *testing.T) {

	box := NewBoundingBoxF(-1, -4, -5, 9, 6, 5)
	left, right := SplitBounds(box)
	assert.Equal(t, NewPoint(-1, -4, -5), left.Min)
	assert.Equal(t, NewPoint(4, 6, 5), left.Max)
	assert.Equal(t, NewPoint(4, -4, -5), right.Min)
	assert.Equal(t, NewPoint(9, 6, 5), right.Max)
}

func TestSplitXWideBoundingBox(t *testing.T) {

	box := NewBoundingBoxF(-1, -2, -3, 9, 5.5, 3)
	left, right := SplitBounds(box)
	assert.Equal(t, NewPoint(-1, -2, -3), left.Min)
	assert.Equal(t, NewPoint(4, 5.5, 3), left.Max)
	assert.Equal(t, NewPoint(4, -2, -3), right.Min)
	assert.Equal(t, NewPoint(9, 5.5, 3), right.Max)
}

func TestSplitYWideBoundingBox(t *testing.T) {

	box := NewBoundingBoxF(-1, -2, -3, 5, 8, 3)
	left, right := SplitBounds(box)
	assert.Equal(t, NewPoint(-1, -2, -3), left.Min)
	assert.Equal(t, NewPoint(5, 3, 3), left.Max)
	assert.Equal(t, NewPoint(-1, 3, -3), right.Min)
	assert.Equal(t, NewPoint(5, 8, 3), right.Max)
}
func TestSplitZWideBoundingBox(t *testing.T) {

	box := NewBoundingBoxF(-1, -2, -3, 5, 3, 7)
	left, right := SplitBounds(box)
	assert.Equal(t, NewPoint(-1, -2, -3), left.Min)
	assert.Equal(t, NewPoint(5, 3, 2), left.Max)
	assert.Equal(t, NewPoint(-1, -2, 2), right.Min)
	assert.Equal(t, NewPoint(5, 3, 7), right.Max)
}

func TestPartitionChildrenOfGroup(t *testing.T) {

	s1 := NewSphere()
	s1.SetTransform(Translate(-2, 0, 0))
	s2 := NewSphere()
	s2.SetTransform(Translate(2, 0, 0))
	s3 := NewSphere()

	g := NewGroup()
	g.AddChild(s1)
	g.AddChild(s2)
	g.AddChild(s3)
	g.Bounds()

	left, right := PartitionChildren(g)
	assert.Equal(t, s1.ID(), left.Children[0].ID())
	assert.Equal(t, s2.ID(), right.Children[0].ID())
	assert.Equal(t, s3.ID(), g.Children[0].ID())

	assert.Equal(t, 1, len(left.Children))
	assert.Equal(t, 1, len(right.Children))
	assert.Equal(t, 1, len(g.Children))
}

func TestCreateSubGroupFromListOfChildren(t *testing.T) {
	s1 := NewSphere()
	s2 := NewSphere()
	g := NewGroup()
	MakeSubGroup(g, s1, s2)
	assert.Len(t, g.Children, 1)
	subgroup := g.Children[0].(*Group)
	assert.True(t, subgroup.Children[0].ID() == s1.ID())
	assert.True(t, subgroup.Children[1].ID() == s2.ID())
}

func TestDividePrimitiveDoesNothing(t *testing.T) {
	/*
		Subdividing a primitive does nothing
		  Given shape ← sphere()
		  When divide(shape, 1)
		  Then shape is a sphere
	*/
	s := NewSphere()
	Divide(s, 1)
	assert.IsType(t, NewSphere(), s)
}

func TestSubdivideGroupPartitionsItsChildren(t *testing.T) {
	/*

	  Then g[0] = s3
	    And subgroup ← g[1]
	    And subgroup is a group
	    And subgroup.count = 2
	    And subgroup[0] is a group of [s1]
	    And subgroup[1] is a group of [s2]
	*/
	s1 := NewSphere()
	s1.SetTransform(Translate(-2, -2, 0))
	s2 := NewSphere()
	s2.SetTransform(Translate(-2, 2, 0))
	s3 := NewSphere()
	s3.SetTransform(Scale(4, 4, 4))

	g := NewGroup()
	g.AddChild(s1)
	g.AddChild(s2)
	g.AddChild(s3)

	Divide(g, 1)

	assert.True(t, g.Children[0].ID() == s3.ID())

	subgroup := g.Children[1].(*Group)
	assert.Equal(t, 2, len(subgroup.Children))
	assert.Equal(t, s1.ID(), subgroup.Children[0].(*Group).Children[0].ID())
	assert.Equal(t, s2.ID(), subgroup.Children[1].(*Group).Children[0].ID())
}

func TestName(t *testing.T) {

	s1 := NewSphere()
	s1.SetTransform(Translate(-2, 0, 0))
	s2 := NewSphere()
	s2.SetTransform(Translate(2, 1, 0))
	s3 := NewSphere()
	s3.SetTransform(Translate(2, -1, 0))
	subgr := NewGroup()
	subgr.AddChildren(s1, s2, s3)

	s4 := NewSphere()

	g := NewGroup()
	g.AddChildren(subgr, s4)

	Divide(g, 3)

	child1 := g.Children[0].(*Group)
	assert.Equal(t, 2, len(child1.Children))
	assert.True(t, subgr == g.Children[0])
	assert.True(t, child1.Children[0].(*Group).Children[0].ID() == s1.ID())
	assert.True(t, child1.Children[1].(*Group).Children[0].ID() == s2.ID())
	assert.True(t, child1.Children[1].(*Group).Children[1].ID() == s3.ID())
	assert.True(t, s4.ID() == g.Children[1].(*Sphere).ID())

}

func TestSubdivideCSGShape(t *testing.T) {
	s1 := NewSphere()
	s1.SetTransform(Translate(-1.5, 0, 0))
	s2 := NewSphere()
	s2.SetTransform(Translate(1.5, 0, 0))
	s3 := NewSphere()
	s3.SetTransform(Translate(0, 0, -1.5))
	s4 := NewSphere()
	s4.SetTransform(Translate(0, 0, 1.5))

	// groups
	left := NewGroup()
	left.AddChildren(s1, s2)

	right := NewGroup()
	right.AddChildren(s3, s4)

	csg := NewCSG("difference", left, right)
	Divide(csg, 1)

	assert.True(t, left.Children[0].(*Group).Children[0].ID() == s1.ID())
	assert.True(t, left.Children[1].(*Group).Children[0].ID() == s2.ID())
	assert.True(t, right.Children[0].(*Group).Children[0].ID() == s3.ID())
	assert.True(t, right.Children[1].(*Group).Children[0].ID() == s4.ID())
}
