package mat

import (
	"sort"
)

func NewRay(origin Tuple4, direction Tuple4) Ray {
	return Ray{Origin: origin, Direction: direction}
}

type Ray struct {
	Origin    Tuple4
	Direction Tuple4
}

// Position multiplies direction of ray with the passed distance and adds the result onto the origin.
// Used for finding the position along a ray.
func Position(r Ray, distance float64) Tuple4 {
	add := MultiplyByScalar(r.Direction, distance)
	pos := Add(r.Origin, add)
	return pos
}

func IntersectRayWithShape(s Shape, r2 Ray) []Intersection {

	// transform ray with inverse of shape transformation matrix to be able to intersect a translated/rotated/skewed shape
	r := TransformRay(r2, Inverse(s.GetTransform()))

	// Call the intersect function provided by the shape implementation (e.g. Sphere, Plane osv)
	return s.IntersectLocal(r)
}

func Hit(intersections []Intersection) (Intersection, bool) {

	// Filter out all negatives
	xs := make([]Intersection, 0)
	for _, i := range intersections {
		if i.T > 0.0 {
			xs = append(xs, i)
		}
	}

	if len(xs) == 0 {
		return Intersection{}, false
	}
	if len(xs) == 1 {
		return xs[0], true
	}
	sort.Slice(xs, func(i, j int) bool {
		return xs[i].T < xs[j].T
	})
	return xs[0], true
}

func TransformRay(r Ray, m1 Mat4x4) Ray {
	origin := MultiplyByTuple(m1, r.Origin)
	direction := MultiplyByTuple(m1, r.Direction)
	return NewRay(origin, direction)
}

func ColorAt(w World, r Ray, remaining int) Tuple4 {
	xs := IntersectWithWorld(w, r)
	if len(xs) > 0 {
		comps := PrepareComputationForIntersection(xs[0], r)
		return ShadeHit(w, comps, remaining)
	} else {
		return black
	}
}

func ShadeHit(w World, comps Computation, remaining int) Tuple4 {
	inShadow := PointInShadow(w, comps.OverPoint)
	color := Lighting(comps.Object.GetMaterial(), comps.Object, w.Light, comps.Point, comps.EyeVec, comps.NormalVec, inShadow)
	reflectedColor := ReflectedColor(w, comps, remaining)
	return Add(color, reflectedColor)
}

func PointInShadow(w World, p Tuple4) bool {
	vecToLight := Sub(w.Light.Position, p)
	distance := Magnitude(vecToLight)

	//movedPos := Add(p, MultiplyByScalar(comps.NormalVec, 0.1))
	ray := NewRay(p, Normalize(vecToLight))
	xs := IntersectWithWorld(w, ray)
	if len(xs) > 0 {
		for _, x := range xs {
			if x.T < distance {
				return true
			}
		}
	}
	return false
}
