package mat

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
