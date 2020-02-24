package mat

import (
	"github.com/eriklupander/rt/internal/pkg/calcstats"
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
	xsCache   Intersections
	bb        *BoundingBox

	CastShadow bool
}

func NewGroup() *Group {
	m1 := New4x4() //NewMat4x4(make([]float64, 16))
	//copy(m1.Elems, IdentityMatrix.Elems)
	inv := New4x4() //NewMat4x4(make([]float64, 16))
	//copy(inv.Elems, IdentityMatrix.Elems)

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

		bb:         NewEmptyBoundingBox(),
		CastShadow: true,
	}
}

func (g *Group) ID() int64 {
	return g.Id
}

func (g *Group) GetTransform() Mat4x4 {
	return g.Transform
}

func (g *Group) GetInverse() Mat4x4 {
	return g.Inverse
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

	if g.bb != nil && !IntersectRayWithBox(ray, g.bb) {
		calcstats.Incr()
		return nil
	}
	//TransformRayPtr(ray, g.Inverse, &g.savedRay)

	g.xsCache = g.xsCache[:0]
	for idx := range g.Children {
		TransformRayPtr(ray, g.Children[idx].GetInverse(), &g.innerRays[idx])
		lxs := g.Children[idx].IntersectLocal(g.innerRays[idx])
		if len(lxs) > 0 {
			g.xsCache = append(g.xsCache, lxs...)
		}
	}

	if len(g.xsCache) > 1 {
		sort.Sort(g.xsCache)
	}
	return g.xsCache
}

func (g *Group) NormalAtLocal(point Tuple4, intersection *Intersection) Tuple4 {
	panic("not applicable to a group")
}

func (g *Group) GetLocalRay() Ray {
	panic("not applicable to a group")
}

func (g *Group) AddChildren(shapes ...Shape) {
	for i := 0; i < len(shapes); i++ {
		g.AddChild(shapes[i])
	}
}

func (g *Group) AddChild(s Shape) {
	g.Children = append(g.Children, s)
	s.SetParent(g)

	// allocate memory for inner rays each time a child is added.
	g.innerRays = append(g.innerRays, NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 0)))

	// recalculate bounds
	g.bb.MergeWith(BoundsOf(s))
}

func (g *Group) Bounds() {
	g.bb = BoundsOf(g)
	//g.bb = TransformBoundingBox(g.bb, g.GetTransform()) // transform by the group's own transform too
}

func (g *Group) CastsShadow() bool {
	return g.CastShadow
}

func (g *Group) GetParent() Shape {
	return g.Parent
}
func (g *Group) SetParent(shape Shape) {
	g.Parent = shape
}
func (g *Group) BoundsToCube() *Cube {
	//v := Sub(g.bb.Max, g.bb.Min)
	//length := Magnitude(v)
	//centre := Position(NewRay(g.bb.Max, v), length / 2)
	TransformBoundingBox(g.bb, g.Transform)
	xscale := (g.bb.Max[0] - g.bb.Min[0]) / 2
	yscale := (g.bb.Max[1] - g.bb.Min[1]) / 2
	zscale := (g.bb.Max[2] - g.bb.Min[2]) / 2
	x := g.bb.Min[0] + xscale
	y := g.bb.Min[1] + yscale
	z := g.bb.Min[2] + zscale

	c := NewCube()
	//c.SetTransform(Translate(-1, 0.25, -1)) //g.bb.Max[0] - g.bb.Min))
	c.SetTransform(g.Transform)
	c.SetTransform(Translate(x, y, z))
	c.SetTransform(Scale(xscale, yscale, zscale))
	m := NewDefaultMaterial()
	m.Transparency = 0.8
	m.Color = NewColor(0.8, 0.7, 0.9)
	c.SetMaterial(m)
	return c
}
