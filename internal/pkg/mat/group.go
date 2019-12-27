package mat

import (
	"math/rand"
	"sort"
)

type Group struct {
	Id        int64
	Transform Mat4x4
	Inverse   Mat4x4
	Material  Material
	Label     string
	Parent    Shape
	Children  []Shape
	savedRay  Ray

	innerRays []Ray
	xsCache   []Intersection
	bb        *BoundingBox
	c         *Cube
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
	inv := NewMat4x4(make([]float64, 16))
	copy(inv.Elems, IdentityMatrix.Elems)

	cachedXs := make([]Intersection, 16)
	innerRays := make([]Ray, 0)

	return &Group{
		Id:        rand.Int63(),
		Transform: m1,
		Inverse:   inv,
		Children:  make([]Shape, 0),
		savedRay:  NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 0)),

		xsCache:   cachedXs,
		innerRays: innerRays,
	}
}

func (g *Group) ID() int64 {
	return g.Id
}

func (g *Group) GetTransform() Mat4x4 {
	return g.Transform
}

func (g *Group) GetInverse() Mat4x4 {
	return g.Inverse //Inverse(g.Transform)
}

func (g *Group) SetTransform(transform Mat4x4) {
	g.Transform = Multiply(g.Transform, transform)
	g.Inverse = Inverse(g.Transform)
}

func (g *Group) GetMaterial() Material {
	panic("not applicable to group")
}

func (g *Group) SetMaterial(material Material) {
	for _, c := range g.Children {
		c.SetMaterial(material)
	}
}

func (g *Group) IntersectLocal(ray Ray) []Intersection {

	TransformRayPtr(ray, g.Inverse, &g.savedRay)

	//xs := make([]Intersection, 0)
	g.xsCache = g.xsCache[:0]
	for idx := range g.Children {
		TransformRayPtr(g.savedRay, g.Children[idx].GetInverse(), &g.innerRays[idx])
		lxs := g.Children[idx].IntersectLocal(g.innerRays[idx])
		if len(lxs) > 0 {
			g.xsCache = append(g.xsCache, lxs...)
		}
	}

	sort.Slice(g.xsCache, func(i, j int) bool {
		return g.xsCache[i].T < g.xsCache[j].T
	})
	return g.xsCache
}

func (g *Group) NormalAtLocal(point Tuple4, intersection *Intersection) Tuple4 {
	panic("not applicable to a group")
}

func (g *Group) GetLocalRay() Ray {
	panic("not applicable to a group")
}

func (g *Group) AddChild(s Shape) {
	g.Children = append(g.Children, s)
	s.SetParent(g)

	// allocate memory for inner rays each time a child is added.
	g.innerRays = append(g.innerRays, NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 0)))
}
