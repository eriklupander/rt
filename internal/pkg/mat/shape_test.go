package mat

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestBaseShapeTransform(t *testing.T) {
	s := NewSphere()
	s.SetTransform(Multiply(s.Transform, Translate(2, 3, 4)))
	assert.Equal(t, Translate(2, 3, 4), s.Transform)
}
func TestBaseShapeMaterial(t *testing.T) {
	s := NewSphere()
	s.SetMaterial(NewMaterial(NewColor(1, 1, 1), 1.0, 0.1, 0.1, 0.1))
	assert.Equal(t, 1.0, s.Material.Ambient)
}
func TestConvertPointFromWorldToObjectSpace(t *testing.T) {
	gr1 := NewGroup()
	gr1.SetTransform(RotateY(math.Pi / 2))
	gr2 := NewGroup()
	gr2.SetTransform(Scale(2, 2, 2))

	gr1.AddChild(gr2)

	sphere := NewSphere()
	sphere.SetTransform(Translate(5, 0, 0))
	gr2.AddChild(sphere)

	p := WorldToObject(sphere, NewPoint(-2, 0, -10))
	assert.Equal(t, 0.0, p.Get(0))
	assert.Equal(t, 0.0, p.Get(1))
	assert.InEpsilon(t, -1.0, p.Get(2), Epsilon)
}

func TestConvertNormalFromObjectToWorldSpace(t *testing.T) {
	gr := NewGroup()
	gr.SetTransform(RotateY(math.Pi / 2))

	gr2 := NewGroup()
	gr2.SetTransform(Scale(1, 2, 3))
	gr.AddChild(gr2)

	sphere := NewSphere()
	sphere.SetTransform(Translate(5, 0, 0))
	gr2.AddChild(sphere)

	n := NormalToWorld(sphere, NewVector(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3))
	assert.InEpsilon(t, 0.2857, n.Get(0), Epsilon)
	assert.InEpsilon(t, 0.4286, n.Get(1), Epsilon)
	assert.InEpsilon(t, -0.8571, n.Get(2), Epsilon)
}

func TestNormalOnChildObject(t *testing.T) {
	gr := NewGroup()
	gr.SetTransform(RotateY(math.Pi / 2))

	gr2 := NewGroup()
	gr2.SetTransform(Scale(1, 2, 3))
	gr.AddChild(gr2)

	sphere := NewSphere()
	sphere.SetTransform(Translate(5, 0, 0))
	gr2.AddChild(sphere)

	n := NormalAt(sphere, NewPoint(1.7321, 1.1547, -5.5774), nil)
	assert.InEpsilon(t, 0.2857, n.Get(0), Epsilon)
	assert.InEpsilon(t, 0.4286, n.Get(1), Epsilon)
	assert.InEpsilon(t, -0.8571, n.Get(2), Epsilon)
	/*
		 normal on a child object
		Given g1 ← group()
		And set_transform(g1, rotation_y(π/2))
		And g2 ← group()
		And set_transform(g2, scaling(1, 2, 3))
		And add_child(g1, g2)
		And s ← sphere()
		And set_transform(s, translation(5, 0, 0))
		And add_child(g2, s)
		When n ← normal_at(s, point(1.7321, 1.1547, -5.5774))
		Then n = vector(0.2857, 0.4286, -0.8571)
	*/
}
