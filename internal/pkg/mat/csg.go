package mat

import (
	"math/rand"
	"sort"
)

func NewCSG(operation string, left, right Shape) *CSG {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	inv := NewMat4x4(make([]float64, 16))
	copy(inv.Elems, IdentityMatrix.Elems)
	c := &CSG{Id: rand.Int63(), Transform: m1, Inverse: inv, Left: left, Right: right, Operation: operation}
	left.SetParent(c)
	right.SetParent(c)
	return c
}

type CSG struct {
	Id        int64
	Transform Mat4x4
	Inverse   Mat4x4
	Left      Shape
	Right     Shape
	Operation string
	Parent    Shape
	Material  Material
}

func (c *CSG) ID() int64 {
	return c.Id
}

func (c *CSG) GetTransform() Mat4x4 {
	return c.Transform
}

func (c *CSG) GetInverse() Mat4x4 {
	return c.Inverse
}

func (c *CSG) SetTransform(transform Mat4x4) {
	c.Transform = Multiply(c.Transform, transform)
	c.Inverse = Inverse(c.Transform)
}

func (c *CSG) GetMaterial() Material {
	return c.Material
}

func (c *CSG) SetMaterial(material Material) {
	c.Material = material
}

func (c *CSG) IntersectLocal(ray Ray) []Intersection {
	leftXs := IntersectRayWithShape(c.Left, ray)
	rightXs := IntersectRayWithShape(c.Right, ray)
	xs := append(leftXs, rightXs...)
	sort.Slice(xs, func(i, j int) bool {
		return xs[i].T < xs[j].T
	})
	return FilterIntersections(c, xs)
}

func (c *CSG) NormalAtLocal(point Tuple4, intersection *Intersection) Tuple4 {
	return Tuple4{}
}

func (c *CSG) GetLocalRay() Ray {
	panic("impl me")
}

func (c *CSG) GetParent() Shape {
	return c.Parent
}

func (c *CSG) SetParent(shape Shape) {
	c.Parent = shape
}
