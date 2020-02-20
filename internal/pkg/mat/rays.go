package mat

import (
	"math"
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

func PositionPtr(r Ray, distance float64, out *Tuple4) {
	add := MultiplyByScalar(r.Direction, distance)
	AddPtr(r.Origin, add, out)
}

func IntersectRayWithShape(s Shape, r2 Ray) []Intersection {

	// transform ray with inverse of shape transformation matrix to be able to intersect a translated/rotated/skewed shape
	r := TransformRay(r2, s.GetInverse())

	// Call the intersect function provided by the shape implementation (e.g. Sphere, Plane osv)
	return s.IntersectLocal(r)
}

func IntersectRayWithShapePtr(s Shape, r2 Ray, in *Ray) []Intersection {

	// transform ray with inverse of shape transformation matrix to be able to intersect a translated/rotated/skewed shape
	TransformRayPtr(r2, s.GetInverse(), in)

	// Call the intersect function provided by the shape implementation (e.g. Sphere, Plane osv)
	return s.IntersectLocal(*in)
}

// Hit finds the first intersection with a positive T (the passed intersections are assumed to have been sorted already)
func Hit(intersections []Intersection) (Intersection, bool) {

	// Filter out all negatives
	//xs := make([]Intersection, 0)
	for i := 0; i < len(intersections); i++ {
		if intersections[i].T > 0.0 {
			return intersections[i], true
			//xs = append(xs, i)
		}
	}

	//if len(xs) == 0 {
	return Intersection{}, false
	//}
	//if len(xs) == 1 {
	//	return xs[0], true
	//}

	//return xs[0], true
}

func TransformRay(r Ray, m1 Mat4x4) Ray {
	origin := MultiplyByTuple(m1, r.Origin)
	direction := MultiplyByTuple(m1, r.Direction)
	return NewRay(origin, direction)
}

func TransformRayPtr(r Ray, m1 Mat4x4, out *Ray) {
	MultiplyByTuplePtr(m1, r.Origin, &out.Origin)
	MultiplyByTuplePtr(m1, r.Direction, &out.Direction)
}

func Schlick(comps Computation) float64 {
	// find the cosine of the angle between the eye and normal vectors using Dot
	cos := Dot(comps.EyeVec, comps.NormalVec)
	// total internal reflection can only occur if n1 > n2
	if comps.N1 > comps.N2 {
		n := comps.N1 / comps.N2
		sin2Theta := (n * n) * (1.0 - (cos * cos))
		if sin2Theta > 1.0 {
			return 1.0
		}
		// compute cosine of theta_t using trig identity
		cosTheta := math.Sqrt(1.0 - sin2Theta)

		// when n1 > n2, use cos(theta_t) instead
		cos = cosTheta
	}
	temp := (comps.N1 - comps.N2) / (comps.N1 + comps.N2)
	r0 := temp * temp
	return r0 + (1-r0)*math.Pow(1-cos, 5)
}
