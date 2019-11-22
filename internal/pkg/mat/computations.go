package mat

func PrepareComputationForIntersection(i Intersection, r Ray) Computation {
	pos := Position(r, i.T)
	eyeVec := Negate(r.Direction)
	normalVec := NormalAtPoint(i.S, pos)
	inside := false
	if Dot(eyeVec, normalVec) < 0 {
		inside = true
		normalVec = Negate(normalVec)
	}
	overPoint := Add(pos, MultiplyByScalar(normalVec, Epsilon))
	return Computation{
		T:         i.T,
		Object:    i.S,
		Point:     pos,
		EyeVec:    eyeVec,
		NormalVec: normalVec,
		Inside:    inside,
		OverPoint: overPoint,
	}

}

type Computation struct {
	T         float64
	Object    Sphere
	Point     Tuple4
	EyeVec    Tuple4
	NormalVec Tuple4
	Inside    bool
	OverPoint Tuple4
}
