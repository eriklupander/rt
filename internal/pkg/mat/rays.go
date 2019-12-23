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

func PositionPtr(r Ray, distance float64, pos *Tuple4) {
	add := MultiplyByScalar(r.Direction, distance)
	AddPtr(r.Origin, add, pos)
}

//var currentRay = Ray{
//	Origin:    NewPoint(0, 0, 0),
//	Direction: NewVector(0, 0, 0),
//}

func IntersectRayWithShape(s Shape, r2 Ray) []Intersection {

	// transform ray with inverse of shape transformation matrix to be able to intersect a translated/rotated/skewed shape
	r := TransformRay(r2, s.GetInverse())
	//copy(currentRay.Origin.Elems, r2.Origin.Elems)
	//MultiplyByTuplePtr(s.GetInverse(), &currentRay.Origin)
	//copy(currentRay.Direction.Elems, r2.Direction.Elems)
	//MultiplyByTuplePtr(s.GetInverse(), &currentRay.Direction)

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

func ColorAt(w World, r Ray, remaining1, remaining2 int) Tuple4 {
	xs := IntersectWithWorld(w, r)
	if len(xs) > 0 {
		comps := PrepareComputationForIntersection(xs[0], r)
		return ShadeHit(w, comps, remaining1, remaining2)
	} else {
		return black
	}
}

func ShadeHit(w World, comps Computation, remaining1, remaining2 int) Tuple4 {
	var surfaceColor = NewColor(0, 0, 0)
	for _, light := range w.Light {
		inShadow := PointInShadow(w, light, comps.OverPoint)
		color := light.Lighting(comps.Object.GetMaterial(), comps.Object, comps.Point, comps.EyeVec, comps.NormalVec, inShadow)
		surfaceColor = Add(surfaceColor, color)
	}
	reflectedColor := ReflectedColor(w, comps, remaining1, remaining2)
	refractedColor := RefractedColor(w, comps, remaining2)

	mat := comps.Object.GetMaterial()
	if mat.Reflectivity > 0.0 && mat.Transparency > 0.0 {
		reflectance := Schlick(comps)
		return Add(Add(surfaceColor, MultiplyByScalar(reflectedColor, reflectance)), MultiplyByScalar(refractedColor, 1-reflectance))
	} else {
		return Add(surfaceColor, Add(reflectedColor, refractedColor))
	}
}

func PointInShadow(w World, light Light, p Tuple4) bool {

	vecToLight := Sub(light.Position, p)
	distance := Magnitude(vecToLight)

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

func Schlick(comps Computation) float64 {
	// find the cosine of the angle between the eye and normal vectors using Dot
	cos := Dot(comps.EyeVec, comps.NormalVec)
	// total internal reflection can only occur if n1 > n2
	if comps.N1 > comps.N2 {
		n := comps.N1 / comps.N2
		sin2Theta := n * n * (1.0 - cos*cos)
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
