package mat

import "math/rand"

func NewCSG(operation string, left, right Shape) *CSG {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	c := &CSG{Id: rand.Int63(), Transform: m1, Left: left, Right: right, Operation: operation}
	left.SetParent(c)
	right.SetParent(c)
	return c
}

type CSG struct {
	Id        int64
	Transform Mat4x4
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

func (c *CSG) SetTransform(transform Mat4x4) {

}

func (c *CSG) GetMaterial() Material {
	return c.Material
}

func (c *CSG) SetMaterial(material Material) {
	c.Material = material
}

func (c *CSG) IntersectLocal(ray Ray) []Intersection {
	return nil
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
