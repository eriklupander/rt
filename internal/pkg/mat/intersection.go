package mat

import (
	"math"
	"sort"
)

type Intersection struct {
	T float64
	S Shape
	U float64
	V float64
}

type Intersections []Intersection

func (xs Intersections) Len() int           { return len(xs) }
func (xs Intersections) Less(i, j int) bool { return xs[i].T < xs[j].T }
func (xs Intersections) Swap(i, j int)      { xs[i], xs[j] = xs[j], xs[i] }

func NewIntersection(t float64, s Shape) Intersection {
	return Intersection{T: t, S: s}
}
func NewIntersectionUV(t float64, s Shape, u, v float64) Intersection {
	return Intersection{T: t, S: s, U: u, V: v}
}

func IntersectionEqual(i1, i2 Intersection) bool {
	return i1.T == i2.T && i1.S.ID() == i2.S.ID()
}

func IntersectWithWorldPtr(w World, r Ray, xs Intersections, inRay *Ray) []Intersection {
	for idx := range w.Objects {
		intersections := IntersectRayWithShapePtr(w.Objects[idx], r, inRay)

		for innerIdx := range intersections {
			xs = append(xs, intersections[innerIdx])
		}
	}
	if len(xs) > 1 {
		sort.Sort(xs)
	}
	return xs
}

func ShadowIntersect(w World, r Ray, distance float64, inRay *Ray) bool {
	for idx := range w.Objects {
		if !w.Objects[idx].CastsShadow() {
			continue
		}
		intersections := IntersectRayWithShapePtr(w.Objects[idx], r, inRay)

		for innerIdx := range intersections {
			if intersections[innerIdx].T > 0.0 && intersections[innerIdx].T < distance {
				return true
			}
		}
	}
	return false
}

func IntersectWithWorldPtrForShadow(w World, r Ray, xs Intersections, inRay *Ray) []Intersection {
	for idx := range w.Objects {
		if !w.Objects[idx].CastsShadow() {
			continue
		}
		intersections := IntersectRayWithShapePtr(w.Objects[idx], r, inRay)

		for innerIdx := range intersections {
			xs = append(xs, intersections[innerIdx])
		}
	}
	sort.Sort(xs)
	return xs
}

func IntersectRayWithBox(ray Ray, BoundingBox *BoundingBox) bool {
	// There is supposed  to be a way to optimize this for fewer checks by looking at early values.
	xtmin, xtmax := checkAxisForBB(ray.Origin.Get(0), ray.Direction.Get(0), BoundingBox.Min[0], BoundingBox.Max[0])
	ytmin, ytmax := checkAxisForBB(ray.Origin.Get(1), ray.Direction.Get(1), BoundingBox.Min[1], BoundingBox.Max[1])
	ztmin, ztmax := checkAxisForBB(ray.Origin.Get(2), ray.Direction.Get(2), BoundingBox.Min[2], BoundingBox.Max[2])

	// Om det största av min-värdena är större än det minsta max-värdet.
	tmin := max(xtmin, ytmin, ztmin)
	tmax := min(xtmax, ytmax, ztmax)
	return tmin < tmax
}
func checkAxisForBB(origin float64, direction float64, minBB, maxBB float64) (min float64, max float64) {
	tminNumerator := minBB - origin
	tmaxNumerator := maxBB - origin
	var tmin, tmax float64
	if math.Abs(direction) >= Epsilon {
		tmin = tminNumerator / direction
		tmax = tmaxNumerator / direction
	} else {
		tmin = tminNumerator * math.Inf(1)
		tmax = tmaxNumerator * math.Inf(1)
	}
	if tmin > tmax {
		// swap
		temp := tmin
		tmin = tmax
		tmax = temp
	}
	return tmin, tmax
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
