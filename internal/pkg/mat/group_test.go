package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGroup(t *testing.T) {
	gr := NewGroup()
	assert.NotNil(t, gr.Transform)
	assert.Equal(t, 0, len(gr.Children))
}

func TestGroup_AddChild(t *testing.T) {
	gr := NewGroup()
	p := NewPlane()
	gr.AddChild(p)

	assert.Equal(t, 1, len(gr.Children))
	assert.Equal(t, p, gr.Children[0])
	assert.Equal(t, gr, p.Parent)
}

func TestIntersectEmptyGroup(t *testing.T) {
	gr := NewGroup()
	ray := NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 1))
	xs := gr.IntersectLocal(ray)
	assert.Len(t, xs, 0)
}

func TestIntersect(t *testing.T) {
	s3 := NewSphere()
	s3.SetTransform(Translate(5, 0, 0))

	ray := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))
	ray = TransformRay(ray, s3.GetInverse())
	x := s3.IntersectLocal(ray)
	assert.Len(t, x, 0)
}

func TestIntersectGroup(t *testing.T) {
	gr := NewGroup()
	s1 := NewSphere()

	s2 := NewSphere()
	s2.SetTransform(Translate(0, 0, -3))

	s3 := NewSphere()
	s3.SetTransform(Translate(5, 0, 0))

	gr.AddChild(s1)
	gr.AddChild(s2)
	gr.AddChild(s3)

	ray := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))

	xs := gr.IntersectLocal(ray)
	assert.Len(t, xs, 4)
	assert.Equal(t, s2.Id, xs[0].S.ID())
	assert.Equal(t, s2.Id, xs[1].S.ID())
	assert.Equal(t, s1.Id, xs[2].S.ID())
	assert.Equal(t, s1.Id, xs[3].S.ID())

}

func TestGroupTransform(t *testing.T) {
	gr := NewGroup()
	gr.SetTransform(Scale(2, 2, 2))
	s := NewSphere()
	s.SetTransform(Translate(5, 0, 0))
	gr.AddChild(s)
	gr.Bounds()
	//gr.bb.Min = NewPoint(8, -2, -2)
	//gr.bb.Max = NewPoint(12, 2, 2)
	r := NewRay(NewPoint(10, 0, -10), NewVector(0, 0, 1))

	xs := gr.IntersectLocal(r)
	assert.Len(t, xs, 2)
	/*
		 Intersecting a transformed group
		Given g ← group()
		And set_transform(g, scaling(2, 2, 2))
		And s ← sphere()
		And set_transform(s, translation(5, 0, 0))
		And add_child(g, s)
		When r ← ray(point(10, 0, -10), vector(0, 0, 1))
		And xs ← intersect(g, r)
		Then xs.count = 2
	*/
}
