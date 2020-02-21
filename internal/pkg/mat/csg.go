package mat

import (
	"github.com/eriklupander/rt/internal/pkg/calcstats"
	"math/rand"
	"sort"
)

func NewCSG(operation string, left, right Shape) *CSG {
	m1 := New4x4()
	inv := New4x4()
	c := &CSG{Id: rand.Int63(),
		Transform:     m1,
		Inverse:       inv,
		Left:          left,
		Right:         right,
		Operation:     operation,
		savedLeftRay:  NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 0)),
		savedRightRay: NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 0)),
		bb:            NewEmptyBoundingBox(),
	}
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

	savedLeftRay  Ray
	savedRightRay Ray

	bb *BoundingBox
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
	if !IntersectRayWithBox(ray, c.bb) {
		calcstats.Incr()
		return nil
	}
	leftXs := IntersectRayWithShapePtr(c.Left, ray, &c.savedLeftRay)
	rightXs := IntersectRayWithShapePtr(c.Right, ray, &c.savedRightRay)
	xs := append(leftXs, rightXs...)
	sort.Sort(Intersections(xs))
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

func (c *CSG) Bounds() {
	c.bb = BoundsOf(c)
}
