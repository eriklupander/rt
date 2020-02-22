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
