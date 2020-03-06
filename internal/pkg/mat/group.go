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
	Mtl       Mtl
	Label     string
	parent    Shape
	Children  []Shape
	savedRay  Ray

	InnerRays   []Ray
	XsCache     Intersections
	BoundingBox *BoundingBox

	CastShadow bool
}

func NewGroup() *Group {
	m1 := New4x4()  //NewMat4x4(make([]float64, 16))
	inv := New4x4() //NewMat4x4(make([]float64, 16))

	cachedXs := make([]Intersection, 16)
	innerRays := make([]Ray, 0)

	return &Group{
		Id:        rand.Int63(),
		Transform: m1,
		Inverse:   inv,
		Children:  make([]Shape, 0),
		savedRay:  NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 0)),

		XsCache:   cachedXs,
		InnerRays: innerRays,

		BoundingBox: NewEmptyBoundingBox(),
		CastShadow:  true,
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

	if g.BoundingBox != nil && !IntersectRayWithBox(ray, g.BoundingBox) {
		calcstats.Incr()
		return nil
	}

	//fmt.Println("testing XS in Group: " + g.Label)
	g.XsCache = g.XsCache[:0]
	for idx := range g.Children {
		TransformRayPtr(ray, g.Children[idx].GetInverse(), &g.InnerRays[idx])
		lxs := g.Children[idx].IntersectLocal(g.InnerRays[idx])
		if len(lxs) > 0 {
			g.XsCache = append(g.XsCache, lxs...)
		}
	}

	if len(g.XsCache) > 1 {
		sort.Sort(g.XsCache)
	}
	return g.XsCache
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
	g.InnerRays = append(g.InnerRays, NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 0)))

	// recalculate bounds
	g.BoundingBox.MergeWith(BoundsOf(s))
}

func (g *Group) Bounds() {
	g.BoundingBox = BoundsOf(g)
}

func (g *Group) CastsShadow() bool {
	return g.CastShadow
}

func (g *Group) GetParent() Shape {
	return g.parent
}
func (g *Group) SetParent(shape Shape) {
	g.parent = shape
}
func (g *Group) BoundsToCube() *Cube {
	TransformBoundingBox(g.BoundingBox, g.Transform)
	xscale := (g.BoundingBox.Max[0] - g.BoundingBox.Min[0]) / 2
	yscale := (g.BoundingBox.Max[1] - g.BoundingBox.Min[1]) / 2
	zscale := (g.BoundingBox.Max[2] - g.BoundingBox.Min[2]) / 2
	x := g.BoundingBox.Min[0] + xscale
	y := g.BoundingBox.Min[1] + yscale
	z := g.BoundingBox.Min[2] + zscale

	c := NewCube()
	c.SetTransform(g.Transform)
	c.SetTransform(Translate(x, y, z))
	c.SetTransform(Scale(xscale, yscale, zscale))

	m := NewDefaultMaterial()
	m.Transparency = 0.8
	m.Color = NewColor(0.8, 0.7, 0.9)
	c.SetMaterial(m)
	return c
}
func (g *Group) Name() string {
	return g.Label
}
