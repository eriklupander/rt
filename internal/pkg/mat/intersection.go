package mat

import (
	"sort"
)

type Intersection struct {
	T float64
	S Shape
	U float64
	V float64
}

func NewIntersection(t float64, s Shape) Intersection {
	return Intersection{T: t, S: s}
}
func NewIntersectionUV(t float64, s Shape, u, v float64) Intersection {
	return Intersection{T: t, S: s, U: u, V: v}
}

func IntersectionEqual(i1, i2 Intersection) bool {
	return i1.T == i2.T && i1.S.ID() == i2.S.ID()
}

// var xs = make([]Intersection, 0)

func IntersectWithWorld(w World, r Ray) []Intersection {
	//xs = xs[:0]
	//xs = nil
	xs := make([]Intersection, 0)
	for idx, _ := range w.Objects {
		intersections := IntersectRayWithShape(w.Objects[idx], r)
		if len(intersections) > 0 {
			for innerIdx := range intersections {
				if intersections[innerIdx].T >= 0.0 {
					xs = append(xs, intersections[innerIdx])
				}
			}
		}
	}
	// Remember that we must sort away negative ones?
	sort.Slice(xs, func(i, j int) bool {
		return xs[i].T < xs[j].T
	})
	return xs
}

func IntersectWithWorldPtr(w World, r Ray, xs []Intersection, inRay *Ray) []Intersection {
	//xs = xs[:0]
	//xs = nil
	for idx, _ := range w.Objects {
		intersections := IntersectRayWithShapePtr(w.Objects[idx], r, inRay)

		for innerIdx := range intersections {
			if intersections[innerIdx].T >= 0.0 {
				xs = append(xs, intersections[innerIdx])
			}
		}

	}
	// Remember that we must sort away negative ones?
	sort.Slice(xs, func(i, j int) bool {
		return xs[i].T < xs[j].T
	})
	return xs
}

func IntersectionAllowed(op string, lhit, inl, inr bool) bool {
	if op == "union" {
		return (lhit && !inr) || (!lhit && !inl)
	}
	if op == "intersection" {
		return (lhit && inr) || (!lhit && inl)
	}
	if op == "difference" {
		return (lhit && !inr) || (!lhit && inl)
	}
	return false
}

func FilterIntersections(csg *CSG, xs []Intersection) []Intersection {
	// begin outside of both children
	inl := false
	inr := false
	// prepare a list to receive the filtered intersections
	result := make([]Intersection, 0)
	for idx, i := range xs {
		// if i.object is part of the "left" child, then lhit is true
		lhit := includes(csg.Left, i.S)
		if IntersectionAllowed(csg.Operation, lhit, inl, inr) {
			result = append(result, xs[idx])
		}
		// depending on which object was hit, toggle either inl or inr
		if lhit {
			inl = !inl
		} else {
			inr = !inr
		}

	}
	return result
}

func includes(left Shape, object Shape) bool {
	switch t := left.(type) {
	case *Group:
		for _, child := range t.Children {
			if child.ID() == object.ID() {
				return true
			} else {
				return includes(child, object)
			}
		}
		return false
	case *CSG:
		a := includes(t.Left, object)
		b := includes(t.Right, object)
		return a || b
	default:
		return left.ID() == object.ID()
	}
}
