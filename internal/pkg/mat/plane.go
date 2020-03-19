package mat

import (
	"math"
	"math/rand"
)

func NewPlane() *Plane {
	m1 := New4x4()  //NewMat4x4(make([]float64, 16))
	inv := New4x4() //NewMat4x4(make([]float64, 16))
	invTranspose := New4x4()
	return &Plane{
		Id:               rand.Int63(),
		Transform:        m1,
		Inverse:          inv,
		InverseTranspose: invTranspose,
		Material:         NewDefaultMaterial(),
		Label:            "Plane",
		savedXs:          make([]Intersection, 1),
		CastShadow:       true,
	}
}

type Plane struct {
	Id               int64
	Transform        Mat4x4
	Inverse          Mat4x4
	InverseTranspose Mat4x4
	Material         Material
	Label            string
	parent           Shape
	savedRay         Ray
	CastShadow       bool

	savedXs []Intersection
}

func (p *Plane) CastsShadow() bool {
	return p.CastShadow
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
func (p *Plane) GetInverseTranspose() Mat4x4 {
	return p.InverseTranspose
}

func (p *Plane) GetMaterial() Material {
	return p.Material
}

// SetTransform passes a pointer to the Plane on which to apply the translation matrix
func (p *Plane) SetTransform(translation Mat4x4) {
	p.Transform = Multiply(p.Transform, translation)
	p.Inverse = Inverse(p.Transform)
	p.InverseTranspose = Transpose(p.Inverse)
}

// SetMaterial passes a pointer to the Plane on which to set the material
func (p *Plane) SetMaterial(m Material) {
	p.Material = m
}

func (p *Plane) IntersectLocal(ray Ray) []Intersection {
	if math.Abs(ray.Direction.Get(1)) < Epsilon {
		return nil
	}
	t := -ray.Origin.Get(1) / ray.Direction.Get(1)
	p.savedXs[0].T = t
	p.savedXs[0].S = p
	return p.savedXs
}

func (p *Plane) NormalAtLocal(point Tuple4, intersection *Intersection) Tuple4 {
	return NewVector(0, 1, 0)
}

func (p *Plane) GetLocalRay() Ray {
	panic("implement me")
}
func (p *Plane) GetParent() Shape {
	return p.parent
}
func (p *Plane) SetParent(shape Shape) {
	p.parent = shape
}
func (p *Plane) Name() string {
	return p.Label
}
