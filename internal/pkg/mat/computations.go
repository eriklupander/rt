package mat

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
