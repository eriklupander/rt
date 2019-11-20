package mat

//func NormalAt(s *Sphere, x, y, z float64) Tuple4 {
//	origin := NewPoint(0, 0, 0)
//	point := NewPoint(x, y, z)
//	return Normalize(*Sub(*point, *origin))
//}
func NormalAtPoint(s Sphere, worldPoint Tuple4) Tuple4 {
	origin := NewPoint(0, 0, 0)

	// transform point on surface of sphere from world to object space.
	objectPoint := MultiplyByTuple(Inverse(s.Transform), worldPoint)

	// vector from the point on the sphere surface to its origin
	objectNormal := Sub(objectPoint, origin)

	// convert normal from object space to world space
	worldNormal := MultiplyByTuple(Transpose(Inverse(s.Transform)), objectNormal)

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
