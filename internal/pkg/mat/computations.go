package mat

func NewComputation() Computation {
	return Computation{
		T:          0,
		Object:     nil,
		Point:      NewPoint(0, 0, 0),
		EyeVec:     NewVector(0, 0, 0),
		NormalVec:  NewVector(0, 0, 0),
		Inside:     false,
		OverPoint:  NewPoint(0, 0, 0),
		UnderPoint: NewPoint(0, 0, 0),
		ReflectVec: NewVector(0, 0, 0),
		N1:         0,
		N2:         0,
	}
}

func PrepareComputationForIntersectionPtr(i Intersection, r Ray, comps *Computation, xs ...Intersection) {
	comps.T = i.T
	comps.Object = i.S
	PositionPtr(r, i.T, &comps.Point)
	NegatePtr(r.Direction, &comps.EyeVec)
	comps.NormalVec = NormalAt(i.S, comps.Point, &i) //  fix
	ReflectPtr(r.Direction, comps.NormalVec, &comps.ReflectVec)

	comps.Inside = false
	if Dot(comps.EyeVec, comps.NormalVec) < 0 {
		comps.Inside = true
		comps.NormalVec = Negate(comps.NormalVec) // fix
	}
	AddPtr(comps.Point, MultiplyByScalar(comps.NormalVec, Epsilon), &comps.OverPoint)
	SubPtr(comps.Point, MultiplyByScalar(comps.NormalVec, Epsilon), &comps.UnderPoint)

	comps.N1 = 1.0
	comps.N2 = 1.0

	var containers = make([]Shape, 0)
	for idx := range xs {
		if xs[idx].S.ID() == i.S.ID() && i.T == xs[idx].T {
			if len(containers) == 0 {
				comps.N1 = 1.0
			} else {
				comps.N1 = containers[len(containers)-1].GetMaterial().RefractiveIndex
			}
		}

		index := indexOf(xs[idx].S, containers)
		if index > -1 {
			copy(containers[index:], containers[index+1:]) // Shift a[i+1:] left one indexs[idx].
			containers[len(containers)-1] = nil            // Erase last element (write zero value).
			containers = containers[:len(containers)-1]    // Truncate slice.
		} else {
			containers = append(containers, xs[idx].S)
		}

		if xs[idx].S.ID() == i.S.ID() && xs[idx].T == i.T {
			if len(containers) == 0 {
				comps.N2 = 1.0
			} else {
				comps.N2 = containers[len(containers)-1].GetMaterial().RefractiveIndex
			}
			break
		}
	}
	//
	//return Computation{
	//	T:          i.T,
	//	Object:     i.S,
	//	Point:      comps.Point,
	//	EyeVec:     eyeVec,
	//	NormalVec:  normalVec,
	//	ReflectVec: reflectVec,
	//	Inside:     inside,
	//	OverPoint:  overPoint,
	//	UnderPoint: underPoint,
	//	N1:         n1,
	//	N2:         n2,
	//}
}

func PrepareComputationForIntersection(i Intersection, r Ray, xs ...Intersection) Computation {
	pos := Position(r, i.T)
	eyeVec := Negate(r.Direction)
	normalVec := NormalAt(i.S, pos, &i)
	reflectVec := Reflect(r.Direction, normalVec)
	inside := false
	if Dot(eyeVec, normalVec) < 0 {
		inside = true
		normalVec = Negate(normalVec)
	}
	overPoint := Add(pos, MultiplyByScalar(normalVec, Epsilon))
	underPoint := Sub(pos, MultiplyByScalar(normalVec, Epsilon))

	n1 := 1.0
	n2 := 1.0

	var containers = make([]Shape, 0)
	for idx := range xs {
		if xs[idx].S.ID() == i.S.ID() && i.T == xs[idx].T {
			if len(containers) == 0 {
				n1 = 1.0
			} else {
				n1 = containers[len(containers)-1].GetMaterial().RefractiveIndex
			}
		}

		index := indexOf(xs[idx].S, containers)
		if index > -1 {
			copy(containers[index:], containers[index+1:]) // Shift a[i+1:] left one indexs[idx].
			containers[len(containers)-1] = nil            // Erase last element (write zero value).
			containers = containers[:len(containers)-1]    // Truncate slice.
		} else {
			containers = append(containers, xs[idx].S)
		}

		if xs[idx].S.ID() == i.S.ID() && xs[idx].T == i.T {
			if len(containers) == 0 {
				n2 = 1.0
			} else {
				n2 = containers[len(containers)-1].GetMaterial().RefractiveIndex
			}
			break
		}
	}

	return Computation{
		T:          i.T,
		Object:     i.S,
		Point:      pos,
		EyeVec:     eyeVec,
		NormalVec:  normalVec,
		ReflectVec: reflectVec,
		Inside:     inside,
		OverPoint:  overPoint,
		UnderPoint: underPoint,
		N1:         n1,
		N2:         n2,
	}
}

//func refractions(i Intersection, xs []Intersection) (float64, float64) {
//	var containers = make([]Shape, 0)
//	var n1, n2 float64
//	for a := 0; a < len(xs); a++	{
//
//		var curr_i = xs[a]
//
//		if curr_i.S.ID() == i.S.ID() && i.T == curr_i.T {
//			if len(containers) == 0 {
//				n1 = 1.0
//			} else {
//				n1 = containers[len(containers)-1].GetMaterial().RefractiveIndex
//			}
//		}
//
//		// if containers includes i.object then
//
//		if indexOfObj := indexOf(curr_i.S, containers); indexOfObj != -1 {
//			copy(containers[indexOfObj:], containers[indexOfObj+1:]) // Shift a[i+1:] left one indexs[idx].
//			containers[len(containers)-1] = nil            // Erase last element (write zero value).
//			containers = containers[:len(containers)-1]    // Truncate slice.
//			//containers.splice(indexOfObj, 1)
//		} else {
//			containers = append(containers, curr_i.S)
//		}
//
//		if curr_i.S.ID() == i.S.ID() && i.T == curr_i.T	{
//
//			if len(containers) == 0 {
//				n2 = 1.0
//			} else {
//				n2 = containers[len(containers)-1].GetMaterial().RefractiveIndex
//			}
//			break
//		}
//	}
//	return n1, n2
//}

func indexOf(s Shape, list []Shape) int {
	for idx := range list {
		if list[idx].ID() == s.ID() {
			return idx
		}
	}
	return -1
}

type Computation struct {
	T          float64
	Object     Shape
	Point      Tuple4
	EyeVec     Tuple4
	NormalVec  Tuple4
	Inside     bool
	OverPoint  Tuple4
	UnderPoint Tuple4
	ReflectVec Tuple4
	N1         float64
	N2         float64
}
