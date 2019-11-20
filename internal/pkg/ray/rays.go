package ray

import (
	"github.com/eriklupander/rt/internal/pkg/mat"
	"math"
	"sort"
)

func New(origin mat.Tuple4, direction mat.Tuple4) Ray {
	return Ray{Origin: origin, Direction: direction}
}

type Ray struct {
	Origin    mat.Tuple4
	Direction mat.Tuple4
}

// Position multiplies direction of ray with the passed distance and adds the result onto the origin.
// Used for finding the position along a ray.
func Position(r Ray, distance float64) mat.Tuple4 {
	add := mat.MultiplyByScalar(r.Direction, distance)
	pos := mat.Add(r.Origin, add)
	return pos
}

func IntersectRayWithSphere(s mat.Sphere, r2 Ray) []Intersection {

	// transform ray with inverse of sphere transformation matrix to be able to intersect a translated/rotated/skewed sphere
	r := Transform(r2, mat.Inverse(s.Transform))

	// this is a vector from the origin of the ray to the center of the sphere at 0,0,0
	sphereToRay := mat.Sub(r.Origin, mat.NewPoint(0, 0, 0))

	// This dot product is
	a := mat.Dot(r.Direction, r.Direction)

	// Take the dot of the direction and the vector from ray origin to sphere center times 2
	b := 2.0 * mat.Dot(r.Direction, sphereToRay)

	// Take the dot of the two sphereToRay vectors and decrease by 1 (is that because the sphere is unit length 1?
	c := mat.Dot(sphereToRay, sphereToRay) - 1.0

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

func Hit(intersections []Intersection) *Intersection {

	// Filter out all negatives
	xs := make([]Intersection, 0)
	for _, i := range intersections {
		if i.T > 0.0 {
			xs = append(xs, i)
		}
	}

	if len(xs) == 0 {
		return nil
	}
	if len(xs) == 1 {
		return &xs[0]
	}
	sort.Slice(xs, func(i, j int) bool {
		return xs[i].T < xs[j].T
	})
	return &xs[0]
}

func Transform(r Ray, m1 mat.Mat4x4) Ray {
	origin := mat.MultiplyByTuple(m1, r.Origin)
	direction := mat.MultiplyByTuple(m1, r.Direction)
	return New(origin, direction)
}
