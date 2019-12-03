package mat

import (
	"math/rand"
	"sort"
)

type Group struct {
	Id        int64
	Transform Mat4x4
	Material  Material
	Label     string
	Parent    Shape
	Children  []Shape
	savedRay  Ray
}

func (g *Group) GetParent() Shape {
	return g.Parent
}
func (g *Group) SetParent(shape Shape) {
	g.Parent = shape
}

func NewGroup() *Group {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	return &Group{Id: rand.Int63(), Transform: m1, Children: make([]Shape, 0)}
}

func (g *Group) ID() int64 {
	return g.Id
}

func (g *Group) GetTransform() Mat4x4 {
	return g.Transform
}

func (g *Group) SetTransform(transform Mat4x4) {
	g.Transform = Multiply(g.Transform, transform)
}

func (g *Group) GetMaterial() Material {
	panic("not applicable to group")
}

func (g *Group) SetMaterial(material Material) {
	panic("not applicable to a group")
}

func (g *Group) IntersectLocal(ray Ray) []Intersection {
	ray = TransformRay(ray, Inverse(g.Transform))

	xs := make([]Intersection, 0)
	for idx := range g.Children {
		ray = TransformRay(ray, Inverse(g.Children[idx].GetTransform()))
		lxs := g.Children[idx].IntersectLocal(ray)
		if len(lxs) > 0 {
			xs = append(xs, lxs...)
		}
	}
	sort.Slice(xs, func(i, j int) bool {
		return xs[i].T < xs[j].T
	})
	return xs
}

func (g *Group) NormalAtLocal(point Tuple4) Tuple4 {
	panic("not applicable to a group")
}

func (g *Group) GetLocalRay() Ray {
	panic("not applicable to a group")
}

func (g *Group) AddChild(s Shape) {
	g.Children = append(g.Children, s)
	s.SetParent(g)
}
