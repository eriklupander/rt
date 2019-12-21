package mat

import (
	"math"
	"math/rand"
)

func NewPlane() *Plane {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	inv := NewMat4x4(make([]float64, 16))
	copy(inv.Elems, IdentityMatrix.Elems)
	return &Plane{Id: rand.Int63(), Transform: m1, Inverse: inv, Material: NewDefaultMaterial(), Label: "Plane"}
}

type Plane struct {
	Id        int64
	Transform Mat4x4
	Inverse   Mat4x4
	Material  Material
	Label     string
	Parent    Shape
	savedRay  Ray
}

func (p *Plane) ID() int64 {
	return p.Id
}
func (p *Plane) GetTransform() Mat4x4 {
	return p.Transform
}
func (p *Plane) GetInverse() Mat4x4 {
	return p.Inverse
}

func (p *Plane) GetMaterial() Material {
	return p.Material
}

// SetTransform passes a pointer to the Plane on which to apply the translation matrix
func (p *Plane) SetTransform(translation Mat4x4) {
	p.Transform = Multiply(p.Transform, translation)
	p.Inverse = Inverse(p.Transform)
}

// SetMaterial passes a pointer to the Plane on which to set the material
func (p *Plane) SetMaterial(m Material) {
	p.Material = m
}

func (p *Plane) IntersectLocal(ray Ray) []Intersection {
	if math.Abs(ray.Direction.Get(1)) < Epsilon {
		return []Intersection{}
	}
	t := -ray.Origin.Get(1) / ray.Direction.Get(1)
	return []Intersection{
		{T: t, S: p},
	}
}

func (p *Plane) NormalAtLocal(point Tuple4, intersection *Intersection) Tuple4 {
	return NewVector(0, 1, 0)
}

func (p *Plane) GetLocalRay() Ray {
	panic("implement me")
}
func (p *Plane) GetParent() Shape {
	return p.Parent
}
func (p *Plane) SetParent(shape Shape) {
	p.Parent = shape
}
