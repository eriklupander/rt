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

	bb *BoundingBox
	c  *Cube
}

//
//func (g *Group) BB() {
//	// recalculate BB after each added child
//	var minX, minY, minZ, maxX, maxY, maxZ float64
//	for _, child := range g.Children {
//		switch tri := child.(type) {
//		case *Triangle:
//			// mins
//			if tri.P1.Get(0) < minX {
//				minX = tri.P1.Get(0)
//			}
//			if tri.P2.Get(0) < minX {
//				minX = tri.P2.Get(0)
//			}
//			if tri.P3.Get(0) < minX {
//				minX = tri.P3.Get(0)
//			}
//			if tri.P1.Get(1) < minY {
//				minY = tri.P1.Get(1)
//			}
//			if tri.P2.Get(1) < minY {
//				minY = tri.P2.Get(1)
//			}
//			if tri.P3.Get(1) < minY {
//				minY = tri.P3.Get(1)
//			}
//			if tri.P1.Get(2) < minZ {
//				minZ = tri.P1.Get(2)
//			}
//			if tri.P2.Get(2) < minZ {
//				minZ = tri.P2.Get(2)
//			}
//			if tri.P3.Get(2) < minZ {
//				minZ = tri.P3.Get(2)
//			}
//
//			// max
//			if tri.P1.Get(0) > maxX {
//				maxX = tri.P1.Get(0)
//			}
//			if tri.P2.Get(0) > maxX {
//				maxX = tri.P2.Get(0)
//			}
//			if tri.P3.Get(0) > maxX {
//				maxX = tri.P3.Get(0)
//			}
//			if tri.P1.Get(1) > maxY {
//				maxY = tri.P1.Get(1)
//			}
//			if tri.P2.Get(1) > maxY {
//				maxY = tri.P2.Get(1)
//			}
//			if tri.P3.Get(1) > maxY {
//				maxY = tri.P3.Get(1)
//			}
//			if tri.P1.Get(2) > maxZ {
//				maxZ = tri.P1.Get(2)
//			}
//			if tri.P2.Get(2) > maxZ {
//				maxZ = tri.P2.Get(2)
//			}
//			if tri.P3.Get(2) > maxZ {
//				maxZ = tri.P3.Get(2)
//			}
//		case *SmoothTriangle:
//
//			// mins
//			if tri.P1.Get(0) < minX {
//				minX = tri.P1.Get(0)
//			}
//			if tri.P2.Get(0) < minX {
//				minX = tri.P2.Get(0)
//			}
//			if tri.P3.Get(0) < minX {
//				minX = tri.P3.Get(0)
//			}
//			if tri.P1.Get(1) < minY {
//				minY = tri.P1.Get(1)
//			}
//			if tri.P2.Get(1) < minY {
//				minY = tri.P2.Get(1)
//			}
//			if tri.P3.Get(1) < minY {
//				minY = tri.P3.Get(1)
//			}
//			if tri.P1.Get(2) < minZ {
//				minZ = tri.P1.Get(2)
//			}
//			if tri.P2.Get(2) < minZ {
//				minZ = tri.P2.Get(2)
//			}
//			if tri.P3.Get(2) < minZ {
//				minZ = tri.P3.Get(2)
//			}
//
//			// max
//			if tri.P1.Get(0) > maxX {
//				maxX = tri.P1.Get(0)
//			}
//			if tri.P2.Get(0) > maxX {
//				maxX = tri.P2.Get(0)
//			}
//			if tri.P3.Get(0) > maxX {
//				maxX = tri.P3.Get(0)
//			}
//			if tri.P1.Get(1) > maxY {
//				maxY = tri.P1.Get(1)
//			}
//			if tri.P2.Get(1) > maxY {
//				maxY = tri.P2.Get(1)
//			}
//			if tri.P3.Get(1) > maxY {
//				maxY = tri.P3.Get(1)
//			}
//			if tri.P1.Get(2) > maxZ {
//				maxZ = tri.P1.Get(2)
//			}
//			if tri.P2.Get(2) > maxZ {
//				maxZ = tri.P2.Get(2)
//			}
//			if tri.P3.Get(2) > maxZ {
//				maxZ = tri.P3.Get(2)
//			}
//		}
//
//	}
//	fmt.Print(minX, minY, minZ, maxX, maxY, maxZ)
//	g.bb = &BoundingBox{
//		PointA: NewPoint(minX, minY, minZ),
//		PointB: NewPoint(maxX, maxY, maxZ),
//	}
//	g.c = NewCube()
//	//g.c.SetTransform(Scale(g.bb.PointB.Get(0) - g.bb.PointA.Get(0), g.bb.PointB.Get(1) - g.bb.PointA.Get(1), g.bb.PointB.Get(2) - g.bb.PointA.Get(2)))
//	g.c.SetTransform(Scale(3.0, 5.1, 3.6))
//}

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

	// allocate space for up to 1024 sub-objects (probably not enough for meshes?)
	innerRays := make([]Ray, 8)
	for i := 0; i < 8; i++ {
		innerRays[i] = NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 0))
	}

	return &Group{
		Id:        rand.Int63(),
		Transform: m1,
		Inverse:   inv,
		Children:  make([]Shape, 0),
		savedRay:  NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 0)),
		innerRays: make([]Ray, 0),
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
	//ray = TransformRay(ray, g.Inverse)
	copy(g.savedRay.Origin.Elems, ray.Origin.Elems)
	MultiplyByTuplePtr(g.GetInverse(), &g.savedRay.Origin)
	copy(g.savedRay.Direction.Elems, ray.Direction.Elems)
	MultiplyByTuplePtr(g.GetInverse(), &g.savedRay.Direction)

	//ray = TransformRay(ray, g.GetInverse())

	// check the bounding box around the group. We should have precomputed this.
	//if g.bb != nil {
	//
	//	lcray := TransformRay(ray, Inverse(g.c.GetTransform()))
	//	if len(g.c.IntersectLocal(lcray)) == 0 {
	//		return []Intersection{}
	//	}
	//}

	xs := make([]Intersection, 0)
	for idx := range g.Children {
		//innerRay := TransformRay(ray, Inverse(g.Children[idx].GetTransform()))
		//innerRay := TransformRay(g.savedRay, g.Children[idx].GetInverse())

		copy(g.innerRays[idx].Origin.Elems, g.savedRay.Origin.Elems)
		MultiplyByTuplePtr(g.Children[idx].GetInverse(), &g.innerRays[idx].Origin)
		copy(g.innerRays[idx].Direction.Elems, g.savedRay.Direction.Elems)
		MultiplyByTuplePtr(g.Children[idx].GetInverse(), &g.innerRays[idx].Direction)

		lxs := g.Children[idx].IntersectLocal(g.innerRays[idx])
		if len(lxs) > 0 {
			xs = append(xs, lxs...)
		}
	}
	sort.Slice(xs, func(i, j int) bool {
		return xs[i].T < xs[j].T
	})
	return xs
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

	// Each time a child is added, also allocate memory for computing the ray transform
	g.innerRays = append(g.innerRays, NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 0)))
}
