package mat

func PrepareComputationForIntersection(i Intersection, r Ray) Computation {
	pos := Position(r, i.T)
	eyeVec := Negate(r.Direction)
	normalVec := NormalAt(i.S, pos)
	reflectVec := Reflect(r.Direction, normalVec)
	inside := false
	if Dot(eyeVec, normalVec) < 0 {
		inside = true
		normalVec = Negate(normalVec)
	}
	overPoint := Add(pos, MultiplyByScalar(normalVec, Epsilon))
	return Computation{
		T:          i.T,
		Object:     i.S,
		Point:      pos,
		EyeVec:     eyeVec,
		NormalVec:  normalVec,
		ReflectVec: reflectVec,
		Inside:     inside,
		OverPoint:  overPoint,
	}

}

type Computation struct {
	T          float64
	Object     Shape
	Point      Tuple4
	EyeVec     Tuple4
	NormalVec  Tuple4
	Inside     bool
	OverPoint  Tuple4
	ReflectVec Tuple4
}
