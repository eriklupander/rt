package mat

import (
	"math"
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

func IntersectRayWithSphere(s Sphere, r2 Ray) []Intersection {

	// transform ray with inverse of sphere transformation matrix to be able to intersect a translated/rotated/skewed sphere
	r := Transform(r2, Inverse(s.Transform))

	// this is a vector from the origin of the ray to the center of the sphere at 0,0,0
	sphereToRay := Sub(r.Origin, NewPoint(0, 0, 0))

	// This dot product is
	a := Dot(r.Direction, r.Direction)

	// Take the dot of the direction and the vector from ray origin to sphere center times 2
	b := 2.0 * Dot(r.Direction, sphereToRay)

	// Take the dot of the two sphereToRay vectors and decrease by 1 (is that because the sphere is unit length 1?
	c := Dot(sphereToRay, sphereToRay) - 1.0

	// calculate the discriminant
	discriminant := (b * b) - 4*a*c
	if discriminant < 0.0 {
		return []Intersection{}
	}

	// finally, find the intersection distances on our ray. Some values:
	//fmt.Printf("c: %v b: %v, a: %v, 2*a: %v sqrt(disc): %v discr: %v\n", c, b, a, 2*a, math.Sqrt(discriminant), discriminant)
	t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2 * a)
	return []Intersection{
		{T: t1, S: s},
		{T: t2, S: s},
	}
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

func Transform(r Ray, m1 Mat4x4) Ray {
	origin := MultiplyByTuple(m1, r.Origin)
	direction := MultiplyByTuple(m1, r.Direction)
	return NewRay(origin, direction)
}

func ShadeHit(w World, comps Computation) Tuple4 {
	inShadow := PointInShadow(w, comps.OverPoint)
	return Lighting(comps.Object.Material, w.Light, comps.Point, comps.EyeVec, comps.NormalVec, inShadow)
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

func ColorAt(w World, r Ray) Tuple4 {
	xs := IntersectWithWorld(w, r)
	if len(xs) > 0 {
		comps := PrepareComputationForIntersection(xs[0], r)
		return ShadeHit(w, comps)
	} else {
		return black
	}
}
