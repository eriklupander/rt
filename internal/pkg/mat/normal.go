package mat

import "math"

func NormalAt(s Shape, worldPoint Tuple4) Tuple4 {

	// transform point on surface of sphere from world to object space.
	localPoint := MultiplyByTuple(Inverse(s.GetTransform()), worldPoint)

	// vector from the point on the sphere surface to its origin
	objectNormal := s.NormalAtLocal(localPoint)

	// convert normal from object space to world space
	worldNormal := MultiplyByTuple(Transpose(Inverse(s.GetTransform())), objectNormal)

	// fix for having a translation in the transform messing up the w part of the world space vector.
	worldNormal.Elems[3] = 0.0
	return Normalize(worldNormal)
}

// in - normal * 2 * dot(in, normal)
func Reflect(vec Tuple4, normal Tuple4) Tuple4 {
	dotScalar := Dot(vec, normal)
	norm := MultiplyByScalar(MultiplyByScalar(normal, 2.0), dotScalar)
	return Sub(vec, norm)
}

func ReflectedColor(w World, comps Computation, remaining1, remaining2 int) Tuple4 {
	if remaining1 == 0 || comps.Object.GetMaterial().Reflectivity == 0.0 {
		return black
	}
	reflectRay := NewRay(comps.OverPoint, comps.ReflectVec)
	remaining1--
	reflectedColor := ColorAt(w, reflectRay, remaining1, remaining2)
	return MultiplyByScalar(reflectedColor, comps.Object.GetMaterial().Reflectivity)
}

func RefractedColor(w World, comps Computation, remaining int) Tuple4 {
	if remaining == 0 || comps.Object.GetMaterial().Transparency == 0.0 {
		return black
	}

	// Find the ratio of first index of refraction to the second.
	nRatio := comps.N1 / comps.N2
	// cos(theta_i) is the same as the dot product of the two vectors
	cosI := Dot(comps.EyeVec, comps.NormalVec)
	// Find sin(theta_t)^2 via trigonometric identity
	sin2Theta := nRatio * nRatio * (1 - cosI*cosI)
	if sin2Theta > 1.0 {
		return black
	}

	// Find cos(theta_t) via trigonometric identity
	cosTheta := math.Sqrt(1.0 - sin2Theta)

	// Compute the direction of the refracted ray
	direction := Sub(MultiplyByScalar(comps.NormalVec, nRatio*cosI-cosTheta), MultiplyByScalar(comps.EyeVec, nRatio))

	// Create the refracted ray
	refractRay := NewRay(comps.UnderPoint, direction)
	// Find the color of the refracted ray, making sure to multiply
	//by the transparency value to account for any opacity
	color := MultiplyByScalar(ColorAt(w, refractRay, remaining-1, remaining-1), comps.Object.GetMaterial().Transparency)

	return color //MultiplyByScalar(color, comps.Object.GetMaterial().RefractiveIndex)
}
